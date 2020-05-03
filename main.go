package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//Date ... type for store date
type Date struct {
	y int
	m int
	d int
}

func main() {
	stockNumber := flag.Int("stock", 6501, "The number of stock")
	start := flag.String("start", "", "Start time in format '20xx.xx'")
	end := flag.String("end", "", "End time in format '20xx.xx'")
	flag.Parse()

	f, err := os.OpenFile("data.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s, err := splitDate(start)
	if err != nil {
		log.Fatalf("Fail to extract start data %s\n", err)
	}
	e, err := splitDate(end)
	if err != nil {
		log.Fatalf("Fail to extract end date %s\n", err)
	}

	tmp := "https://info.finance.yahoo.co.jp/history/?code=%d.T&sy=%d&sm=%d&sd=%d&ey=%d&em=%d&ed=%d&tm=d"

	addr := fmt.Sprintf(tmp,
		*stockNumber,
		s.y, s.m, s.d,
		e.y, e.m, e.d,
	)

	resp, err := http.Get(addr)
	if err != nil {
		log.Fatalf("main : %v\n", err)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("main : %v\n", err)
	}

	defer resp.Body.Close()

	pages := make(chan int, 5)
	go parsePage(doc, pages)
	for page := range pages {
		fmt.Println(page)
	}
}

func splitDate(date *string) (Date, error) {
	var rst Date
	var err error
	dateSlice := strings.Split(*date, ".")

	if len(dateSlice) != 3 {
		return rst, fmt.Errorf("Invalid Format in Date")
	}

	rst.y, err = strconv.Atoi(dateSlice[0])
	if err != nil {
		return rst, fmt.Errorf("Invalid Format on Year")
	}

	rst.m, err = strconv.Atoi(dateSlice[1])
	if err != nil {
		return rst, fmt.Errorf("Invalid Format on Month")
	}

	rst.d, err = strconv.Atoi(dateSlice[2])
	if err != nil {
		return rst, fmt.Errorf("Invalid Format on Month")
	}

	return rst, nil
}

func parsePrice(doc *html.Node) {

	// Parse Path for price
	PriceProperty := []Propety{
		{"class", "padT12"},
		{"data", "table"},
		{"data", "tbody"},
		{"data", "tr"},
	}

	nodes, err := (*DomNode)(doc).Select(PriceProperty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsePrice: find Target err: %v", err)
	}
	stockData := make([]stock, 0)
	for _, node := range nodes {
		datas := node.SelectAll(Propety{"data", "td"})
		if len(datas) > 0 {
			stockData = append(stockData,
				stock{
					"6501",
					parseDate(datas[0].Content()),
					parseStockVal(datas[1].Content()),
					parseStockVal(datas[2].Content()),
					parseStockVal(datas[3].Content()),
					parseStockVal(datas[4].Content()),
					parseStockVal(datas[5].Content()),
					parseStockVal(datas[6].Content())})
		}
	}

	for n, stock := range stockData {
		fmt.Printf("%d\t%s\n", n, &stock)
	}
}

func parsePage(doc *html.Node, ch chan int) {
	PageProperty := []Propety{
		{"class", "ymuiPagingBottom"},
	}

	nodes, err := (*DomNode)(doc).Select(PageProperty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsePage: find Target err: %v", err)
	}

	for _, node := range nodes {
		datas := node.SelectAll(Propety{"data", "a"})
		for _, data := range datas {
			content := data.Content()
			page, err := strconv.Atoi(content)
			if err == nil {
				ch <- page
			}
		}
	}
	close(ch)
}
