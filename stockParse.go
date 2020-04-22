package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//Propety is the property of the element
type Propety struct {
	Type string
	Val  string
}

//DomNode : wrap html.Node structure with searching function
type DomNode html.Node

func writeString(w io.Writer, s string) (n int, err error) {
	type stringWriter interface {
		WriteString(string) (n int, err error)
	}
	if sw, ok := w.(stringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}

//Select : select the DOM from slice of property sturcture
func Select(n *html.Node, targets []Propety) ([]*html.Node, error) {
	htmlNodeStack := findParent(n, targets[0])
	fmt.Println(len(htmlNodeStack))
	targets = targets[1:]
	fmt.Println(targets)

	for _, target := range targets {
		tmpStack := make([]*html.Node, 0)
		for _, htmlNode := range htmlNodeStack {
			tmpStack = append(tmpStack, SelectAll(htmlNode, target)...)
		}
		htmlNodeStack = tmpStack
	}

	return htmlNodeStack, nil
}

//SelectAll : select the elements under html node by property structure
func SelectAll(n *html.Node, target Propety) []*html.Node {
	var rst []*html.Node

	// Travel to other element
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// if find the wanted elements, put into rst
		if target.Type == "data" {
			if c.Data == target.Val {
				rst = append(rst, c)
			}
		} else {
			for _, attr := range c.Attr {
				// TODO: Deal with multiple class
				if attr.Key == target.Type && strings.Contains(attr.Val, target.Val) {
					rst = append(rst, c)
					break
				}
			}
		}
	}
	return rst
}

func findParent(n *html.Node, target Propety) []*html.Node {
	var rst []*html.Node
	attrs := n.Attr
	// if find the wanted elements, put into rst
	if target.Type == "data" {
		if n.Data == target.Val {
			rst = append(rst, n)
		}
	}

	for _, attr := range attrs {
		// TODO: Deal with multiple class
		if attr.Key == target.Type && strings.Contains(attr.Val, target.Val) {
			rst = append(rst, n)
			break
		}
	}

	// Travel to other element
	if n.Type == html.ElementNode || n.Type == html.DocumentNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			tmp := findParent(c, target)
			if len(tmp) > 0 {
				rst = append(rst, tmp...)
			}
		}
	}
	return rst
}

//PrintNode Print the elemets under specific html node
func PrintNode(w io.Writer, n *html.Node, padding string) {
	s := strings.TrimSpace(n.Data)
	attrs := n.Attr
	attrString := ""
	for _, attr := range attrs {
		attrString += fmt.Sprintf("%s: \"%s\" ", attr.Key, attr.Val)
	}
	if len(s) > 0 {
		if n.Type != 1 {
			writeString(w, fmt.Sprintf("%s<%s %s>\n", padding, s, attrString))
		} else {
			writeString(w, padding+s)
		}
	}
	if n.Type == html.ElementNode || n.Type == html.DocumentNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			PrintNode(w, c, padding+"  ")
		}
	}
	if n.Type != 1 && len(s) > 0 {
		writeString(w, fmt.Sprintf("%s</%s>\n", padding, n.Data))
	}
}

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
	nodes, err := Select(doc, TargetProperty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "find Target err: %v\n", err)
	}

	for i, node := range nodes {
		fmt.Printf("Node: %d\n", i)
		PrintNode(os.Stdout, node, "")
	}

}
