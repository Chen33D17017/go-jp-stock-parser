package main

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	"os"
	"regexp"
)

type stock struct {
	name string
	date time.Time
	open float64
	high float64
	low float64
	end float64
	volumn float64
	close float64
}


func (s *stock) String() string {
	return fmt.Sprintf("%s @%s: %f", s.name, s.date, s.close)

}


func parseStockVal(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	rst, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Convert error %s", err)
	}
	return rst
}


func parseDate(s string) time.Time {
		re := regexp.MustCompile(`\d+`)
		reRst := re.FindAllString(s, -1)
		year, yearErr := strconv.Atoi(reRst[0])
		month, monthErr := strconv.Atoi(reRst[1])
		day, dayErr := strconv.Atoi(reRst[2])
		if dayErr != nil || monthErr != nil || yearErr != nil{
			panic("parseDate: fail to convert date")
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		
}