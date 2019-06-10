package der

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// time for cache unit
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
)

// 时间格式化2006-01-02 15:04:05
type JsonTime time.Time

// 实现它的json序列化方法
func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// 反序列化方法 https://stackoverflow.com/questions/45303326/how-to-parse-non-standard-time-format-from-json-in-golang
func (t *JsonTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	ti, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = JsonTime(ti)
	return nil
}

// insert problem https://github.com/jinzhu/gorm/issues/1611#issuecomment-329654638%E3%80%82
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to JsonTime", v)
}

// 时间格式化2006-01-02
type JsonDate time.Time

// 实现它的json序列化方法
func (t JsonDate) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// 反序列化
func (t *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	ti, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*t = JsonDate(ti)
	return nil
}

func (t JsonDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *JsonDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonDate(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to JsonDate", v)
}
