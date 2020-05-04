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

	worklist := make(chan []int, 1)
	unseenPage := make(chan int)
	var n int

	// First page
	worklist <- []int{1}
	n++

	for i := 0; i < 20; i++ {
		go func() {
			for page := range unseenPage {
				targetPage := fmt.Sprintf("%s&p=%d", addr, page)
				doc, err := getDoc(targetPage)
				if err != nil {
					log.Fatalf("Parse page: %d, %v", page, err)
				}
				pages, err := parsePage(doc)
				if err != nil {
					log.Fatalf("Parse page: %d, %v", page, err)
				}
				parsePrice(*stockNumber, doc)
				go func() {
					worklist <- pages
				}()
			}
		}()
	}

	pageSeen := make(map[int]bool)
	for ; n > 0; n-- {
		list := <-worklist
		// for list := range worklist {
		for _, page := range list {
			if !pageSeen[page] {
				pageSeen[page] = true
				n++
				unseenPage <- page
			}
		}
		//}
	}

}

func getDoc(addr string) (*html.Node, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, fmt.Errorf("getDoc: get response %v", err)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getDoc: parse doc %v", err)
	}

	return doc, nil
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
		return rst, fmt.Errorf("Invalid Format on Date")
	}

	return rst, nil
}

func parsePrice(name int, doc *html.Node) {

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
		if len(datas) == 7 {
			stockData = append(stockData,
				stock{
					fmt.Sprintf("%d", name),
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

func parsePage(doc *html.Node) ([]int, error) {
	var rst []int
	PageProperty := []Propety{
		{"class", "ymuiPagingBottom"},
	}

	nodes, err := (*DomNode)(doc).Select(PageProperty)
	if err != nil {
		return nil, fmt.Errorf("parsePage: select Target err: %v", err)
	}

	for _, node := range nodes {
		datas := node.SelectAll(Propety{"data", "a"})
		for _, data := range datas {
			content := data.Content()
			page, err := strconv.Atoi(content)
			if err == nil {
				rst = append(rst, page)
			}
		}
	}
	return rst, nil
}
