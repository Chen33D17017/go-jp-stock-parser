package main

import (
	"fmt"
	"io"
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
func (n *DomNode) Select(targets []Propety) ([]*DomNode, error) {
	htmlNodeStack := n.findParent(targets[0])
	targets = targets[1:]

	for _, target := range targets {
		tmpStack := make([]*DomNode, 0)
		for _, htmlNode := range htmlNodeStack {
			tmpStack = append(tmpStack, htmlNode.SelectAll(target)...)
		}
		htmlNodeStack = tmpStack
	}

	return htmlNodeStack, nil
}

// Content : Return all text under html.node
func (n *DomNode) Content() string {
	rst := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.TextNode {
			rst += (*DomNode)(c).Content()
		} else {
			rst += c.Data
		}
	}
	return rst
}

//SelectAll : select the elements under html node by property structure
func (n *DomNode) SelectAll(target Propety) []*DomNode {
	var rst []*DomNode

	// Travel to other element
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// if find the wanted elements, put into rst
		if target.Type == "data" {
			if c.Data == target.Val {
				rst = append(rst, (*DomNode)(c))
			}
		} else {
			for _, attr := range c.Attr {
				// TODO: Deal with multiple class
				if attr.Key == target.Type && strings.Contains(attr.Val, target.Val) {
					rst = append(rst, (*DomNode)(c))
					break
				}
			}
		}
	}
	return rst
}

func (n *DomNode) findParent(target Propety) []*DomNode {
	var rst []*DomNode
	attrs := n.Attr
	// if find the wanted elements, put into rst
	if target.Type == "data" {
		if n.Data == target.Val {
			rst = append(rst, (*DomNode)(n))
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
			tmp := (*DomNode)(c).findParent(target)
			if len(tmp) > 0 {
				rst = append(rst, tmp...)
			}
		}
	}
	return rst
}

//PrintNode Print the elemets under specific html node
func (n *DomNode) PrintNode(w io.Writer, padding string) {
	s := strings.TrimSpace(n.Data)
	attrs := n.Attr
	attrString := ""
	for _, attr := range attrs {
		attrString += fmt.Sprintf("%s: \"%s\" ", attr.Key, attr.Val)
	}
	if len(s) > 0 {
		if n.Type != html.TextNode {
			writeString(w, fmt.Sprintf("%s<%s %s>\n", padding, s, attrString))
		} else {
			writeString(w, padding+s)
		}
	}
	if n.Type == html.ElementNode || n.Type == html.DocumentNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			(*DomNode)(c).PrintNode(w, padding+"  ")
		}
	}
	if n.Type != html.TextNode && len(s) > 0 {
		writeString(w, fmt.Sprintf("%s</%s>\n", padding, n.Data))
	}
}
