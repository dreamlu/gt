package time

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

const (
	Day        = 24 * time.Minute
	Week       = 7 * Day
	Layout     = "2006-01-02 15:04:05"     // datetime
	LayoutN    = "2006-01-02 15:04:05.000" // datetime(3)
	LayoutDate = "2006-01-02"              // date
	LayoutYM   = "2006-01"                 // date
	LayoutS    = "15:04:05"                // time
)

// CTime china time/date
// format Layout
type CTime time.Time

func (t CTime) MarshalJSON() ([]byte, error) {
	return marshalJSON[CTime](t)
}

func (t *CTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	if len(s) <= 10 {
		s = fmt.Sprintf("%s 00:00:00", s)
	}
	ti, err := parse(Layout, s)
	if err != nil {
		return err
	}
	*t = CTime(ti)
	return nil
}

// Value insert problem https://github.com/go-gorm/gorm/issues/1611#issuecomment-329654638
func (t CTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CTime) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CTime", v)
}

func (CTime) GormDataType() string {
	return "datetime"
}

func (CTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "timestamp" // timestamptz UTC
	default:
		return "datetime" // timestamp UTC
	}
}

func (t CTime) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(Layout)
}

func (t CTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t CTime) Time() time.Time {
	return time.Time(t)
}

// CDate format LayoutDate
type CDate time.Time

func (t CDate) MarshalJSON() ([]byte, error) {
	return marshalJSON[CDate](t)
}

func (t *CDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	ti, err := parse(LayoutDate, s)
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

func (t *CDate) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CDate(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CDate", v)
}

func (CDate) GormDataType() string {
	return "date"
}

func (t CDate) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(LayoutDate)
}

func (t CDate) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t CDate) Time() time.Time {
	return time.Time(t)
}

// CNTime format LayoutN
type CNTime time.Time

func (t CNTime) MarshalJSON() ([]byte, error) {
	return marshalJSON[CNTime](t)
}

func (t *CNTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	ti, err := parse(LayoutN, s)
	if err != nil {
		return err
	}
	*t = CNTime(ti)
	return nil
}

func (t CNTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CNTime) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CNTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CNTime", v)
}

func (CNTime) GormDataType() string {
	return "datetime(3)"
}

func (t CNTime) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(LayoutN)
}

func (t CNTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t CNTime) Time() time.Time {
	return time.Time(t)
}

// CSTime format LayoutS
type CSTime time.Time

func (t CSTime) MarshalJSON() ([]byte, error) {
	return marshalJSON[CSTime](t)
}

func (t *CSTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	ti, err := parse(LayoutS, s)
	if err != nil {
		return err
	}
	*t = CSTime(ti)
	return nil
}

func (t CSTime) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *CSTime) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CSTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CSTime", v)
}

func (CSTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "time"
	default:
		return "time;" // GormDataType gorm bug mysql time to CSTime
	}
}

func (t CSTime) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(LayoutS)
}

func (t CSTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t CSTime) Time() time.Time {
	return time.Time(t)
}

// CYM format LayoutYM
type CYM time.Time

func (t CYM) MarshalJSON() ([]byte, error) {
	return marshalJSON[CYM](t)
}

func (t *CYM) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	ti, err := parse(LayoutYM, s)
	if err != nil {
		return err
	}
	*t = CYM(ti)
	return nil
}

func (t CYM) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CYM) Scan(v any) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CYM(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CYM", v)
}

func (CYM) GormDataType() string {
	return "date"
}

func (t CYM) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(LayoutYM)
}

func (t CYM) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t CYM) Time() time.Time {
	return time.Time(t)
}
