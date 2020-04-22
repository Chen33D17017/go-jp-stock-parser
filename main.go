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

	for _, node := range nodes {
		node.PrintNode(os.Stdout, "")
	}
}