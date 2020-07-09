package time

import (
	"fmt"
	"time"
)

// string to CTime
func ParseCTime(value string) CTime {
	loc, _ := time.LoadLocation("Local")
	ti, err := time.ParseInLocation(Layout, value, loc)
	if err != nil {
		fmt.Println(err)
	}
	return CTime(ti)
}

// string to CDate
func ParseCDate(value string) CDate {
	loc, _ := time.LoadLocation("Local")
	ti, err := time.ParseInLocation(LayoutDate, value, loc)
	if err != nil {
		fmt.Println(err)
	}
	return CDate(ti)
}
