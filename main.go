package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//Date ... type for store date
type Date struct {
	y int
	m int
	d int
}

type config struct {
	StoreType     string `json:"type"`
	ID            string `json:"id"`
	PW            string `json:"password"`
	DB            string `json:"database"`
	DataTable     string `json:"dataTable"`
	NameTable     string `json:"nameTable"`
	CategoryTable string `json:"categoryTable"`
}

func main() {
	startTime := time.Now()
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	stockNumber := flag.Int("stock", 9201, "The number of stock")
	start := flag.String("start", "", "Start time in format '20xx.xx'")
	end := flag.String("end", "", "End time in format '20xx.xx'")
	storageConfig := flag.String("file", "config.json", "Config file for storage, default with config.json")

	flag.Parse()

	// log file
	f, err := os.OpenFile("data.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Read the flag for start & end time
	s, err := splitDate(start)
	if err != nil {
		log.Fatalf("Extract start data error: %s\n", err)
	}
	e, err := splitDate(end)
	if err != nil {
		log.Fatalf("Extract end date error: %s\n", err)
	}

	// Read config file for databas
	file, err := ioutil.ReadFile(*storageConfig)
	if err != nil {
		log.Fatalf("Read File err: %s", err.Error())
	}

	configData := config{}
	err = json.Unmarshal([]byte(file), &configData)
	if err != nil {
		log.Fatalf("Read json err: %s", err.Error())
	}
	// fmt.Println(configData)
	dbm, err := NewDBManager(configData)
	if err != nil {
		log.Fatalf("main: fail to connect db: %s", err)
	}
	defer dbm.Close()

	addr := "https://info.finance.yahoo.co.jp/history/?code=%d.T&sy=%d&sm=%d&sd=%d&ey=%d&em=%d&ed=%d&tm=d"

	addr = fmt.Sprintf(addr,
		*stockNumber,
		s.y, s.m, s.d,
		e.y, e.m, e.d,
	)

	worklist := make(chan []int, 1)
	unseenPage := make(chan int)

	// Parse the start point
	worklist <- []int{1}
	n := 1

	doc, err := getDoc(addr)
	if err != nil {
		log.Fatalf("Parse Page err : %s, %v", addr, err)
	}
	info := parseInfo(*stockNumber, doc)
	// check whether stock in database
	exist, err := dbm.checkStockExist(info, configData.NameTable)
	if err != nil {
		log.Fatalf("main: check stock err: %s", err)
	}
	if !exist {
		dbm.newStock(info, configData)
	}

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
				parsePrice(info, doc, dbm)
				go func() {
					worklist <- pages
				}()
			}
		}()
	}

	pageSeen := make(map[int]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, page := range list {
			if !pageSeen[page] {
				pageSeen[page] = true
				n++
				unseenPage <- page
			}
		}
	}
	fmt.Printf("Runtime: %v", time.Since(startTime))
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

// return name and category (string, string)
func parseInfo(stockID int, doc *html.Node) stockInfo {
	NameSelector := []Selector{{"class", "symbol"}, {"data", "h1"}}
	CatgorySelector := []Selector{{"class", "category"}, {"data", "a"}}

	nameNodes, err := (*DomNode)(doc).Select(NameSelector)
	if err != nil || len(nameNodes) < 1 {
		fmt.Fprintf(os.Stderr, "parseInfo: find Target err: %s", err.Error())
	}

	categoryNodes, err := (*DomNode)(doc).Select(CatgorySelector)
	if err != nil || len(categoryNodes) < 1 {
		fmt.Fprintf(os.Stderr, "parseName: find Target err: %s", err.Error())
	}

	return stockInfo{stockID, nameNodes[0].Content(), categoryNodes[0].Content()}
}

func parsePrice(info stockInfo, doc *html.Node, dbm DBManager) {

	// Parse Path for price
	PriceSelector := []Selector{
		{"class", "padT12"},
		{"data", "table"},
		{"data", "tbody"},
		{"data", "tr"},
	}

	nodes, err := (*DomNode)(doc).Select(PriceSelector)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsePrice: find Target err: %v", err)
	}
	stockData := make([]stock, 0)
	for _, node := range nodes {
		datas := node.SelectAll(Selector{"data", "td"})
		dataSet := [6]float64{}
		var date time.Time
		if len(datas) == 7 {
			for i, data := range datas {
				if i == 0 {
					date = parseDate(data.Content())
				} else {
					dataSet[i-1] = parseStockVal(data.Content())
				}
			}
			stockData = append(stockData, stock{info, date, dataSet})
		}
	}

	for _, stock := range stockData {
		err := dbm.insertStock(stock)
		if err != nil {
			log.Fatalln(err)
		}
		//fmt.Printf("%d\t%s\n", n, &stock)
	}
}

func parsePage(doc *html.Node) ([]int, error) {
	var rst []int
	PageProperty := []Selector{
		{"class", "ymuiPagingBottom"},
	}

	nodes, err := (*DomNode)(doc).Select(PageProperty)
	if err != nil {
		return nil, fmt.Errorf("parsePage: select Target err: %v", err)
	}

	for _, node := range nodes {
		datas := node.SelectAll(Selector{"data", "a"})
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
