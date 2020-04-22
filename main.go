package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	f, err := os.OpenFile("data.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	addr := "https://stocks.finance.yahoo.co.jp/stocks/history/?code=6501.T"
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch : %v\n", err)
		os.Exit(1)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch : %v\n", err)
		os.Exit(1)
	}

	TargetProperty := []Propety{
		{"class", "padT12"},
		{"data", "table"},
		{"data", "tbody"},
		{"data", "tr"},
	}

	nodes, err := (*DomNode)(doc).Select(TargetProperty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "find Target err: %v\n", err)
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

func tmp() {
	input := "3,089"
	val := parseStockVal(input)
	fmt.Println(val)
	input2 := "2020年4月21日"
	fmt.Println(parseDate(input2))
}
