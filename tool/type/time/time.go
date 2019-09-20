package time

import (
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
	"time"
)

// time.Duration expend
const (
	Day  = 24 * time.Minute
	Week = 7 * Day
)

// china time/date
// 时间格式化2006-01-02 15:04:05
type CTime time.Time

// 实现它的json序列化方法
func (t CTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// 反序列化方法 https://stackoverflow.com/questions/45303326/how-to-parse-non-standard-time-format-from-json-in-golang
func (t *CTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	ti, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = CTime(ti)
	return nil
}

// insert problem https://github.com/jinzhu/gorm/issues/1611#issuecomment-329654638%E3%80%82
func (t CTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CTime", v)
}

// must sure MarshalJSON is right
// to string
func (t CTime) String() string {
	// must sure MarshalJSON is right
	b, err := t.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

// 时间格式化2006-01-02
type CDate time.Time

// 实现它的json序列化方法
func (t CDate) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// 反序列化
func (t *CDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	ti, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*t = CDate(ti)
	return nil
}

func (t CDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CDate(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CDate", v)
}

// must sure MarshalJSON is right
// to string
func (t CDate) String() string {
	// must sure MarshalJSON is right
	b, err := t.MarshalJSON()
	if err != nil {
		log.Println(err)
	}
	return string(b)
}
