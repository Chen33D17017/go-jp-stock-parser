package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Date ... type for store date
type Date struct {
	time.Time
}

//NewDate : generate new date struct
func NewDate(y, m, d int) Date {
	return Date{time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)}
}

func (date Date) getTime() time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

//ParseDate : Parse the date from "20xx-xx-xx"
func ParseDate(s string) (Date, error) {
	var err error
	var rst Date
	dateSlice := strings.Split(s, "-")

	if len(dateSlice) != 3 {
		return rst, fmt.Errorf("Invalid Format in Date")
	}

	y, err := strconv.Atoi(dateSlice[0])
	if err != nil {
		return rst, fmt.Errorf("Invalid Format on Year")
	}

	m, err := strconv.Atoi(dateSlice[1])
	if err != nil {
		return rst, fmt.Errorf("Invalid Format on Month")
	}

	d, err := strconv.Atoi(dateSlice[2])
	if err != nil {
		return rst, fmt.Errorf("Invalid Format on Day")
	}

	rst = NewDate(y, m, d)
	return rst, nil
}

func (date Date) String() string {
	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
}
