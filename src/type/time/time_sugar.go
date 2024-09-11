package time

import (
	"fmt"
	"time"
)

// ParseCTime string to CTime
// must Layout
func ParseCTime(value string) CTime {
	ti, err := parse(Layout, value)
	if err != nil {
		fmt.Println(err)
	}
	return CTime(ti)
}

// ParseCDate string to CDate
// must LayoutDate
func ParseCDate(value string) CDate {
	ti, err := parse(LayoutDate, value)
	if err != nil {
		fmt.Println(err)
	}
	return CDate(ti)
}

// ParseCSTime string to CSTime
// must LayoutS
func ParseCSTime(value string) CSTime {
	ti, err := parse(LayoutS, value)
	if err != nil {
		fmt.Println(err)
	}
	return CSTime(ti)
}

func CTimeNow() CTime {
	return CTime(time.Now())
}

func CDateNow() CDate {
	return CDate(time.Now())
}
