package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type stock struct {
	stockInfo
	date    Date
	dataSet [6]float64
}

type stockInfo struct {
	id       int
	name     string
	category string
}

func (s *stock) String() string {
	return fmt.Sprintf("%s %s %d@%s: %f", s.name, s.category, s.id, s.TradeDate(), s.dataSet)
}

func getStockVal(s string) float64 {

	s = strings.ReplaceAll(s, ",", "")
	rst, err := strconv.ParseFloat(s, 64)

	if err != nil {
		fmt.Printf("Convert error with float: %v", err)
		return rst
	}
	return rst
}

func (s *stock) TradeDate() string {
	year, month, day := s.date.Date()
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func getDate(s string) Date {
	re := regexp.MustCompile(`\d+`)
	reRst := re.FindAllString(s, -1)
	year, yearErr := strconv.Atoi(reRst[0])
	month, monthErr := strconv.Atoi(reRst[1])
	day, dayErr := strconv.Atoi(reRst[2])
	if dayErr != nil || monthErr != nil || yearErr != nil {
		panic("parseDate: fail to convert date")
	}
	return NewDate(year, month, day)
}

//func (s *stock)
