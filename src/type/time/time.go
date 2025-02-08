package time

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
	LayoutZero = "0000-00-00 00:00:00.000000000"
)

// CTime china time/date
// format Layout
type CTime time.Time

func (t CTime) MarshalJSON() ([]byte, error) {
	return marshalJSON[CTime](t)
}

func (t *CTime) UnmarshalJSON(b []byte) (err error) {
	*t, err = unmarshalJSON[CTime](Layout, b)
	return
}

func (t *CTime) UnmarshalParam(param string) (err error) {
	*t, err = unmarshalParam[CTime](Layout, param)
	return
}

// Value insert problem https://github.com/go-gorm/gorm/issues/1611#issuecomment-329654638
func (t CTime) Value() (driver.Value, error) {
	return value[CTime](t)
}

func (t *CTime) Scan(v any) (err error) {
	*t, err = scan[CTime](Layout, v)
	return
}

func (CTime) GormDataType() string {
	return "datetime"
}

func (CTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return gormType(db)
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

func (t *CDate) UnmarshalJSON(b []byte) (err error) {
	*t, err = unmarshalJSON[CDate](LayoutDate, b)
	return
}

func (t *CDate) UnmarshalParam(param string) (err error) {
	*t, err = unmarshalParam[CDate](LayoutDate, param)
	return
}

func (t CDate) Value() (driver.Value, error) {
	return value[CDate](t)
}

func (t *CDate) Scan(v any) (err error) {
	*t, err = scan[CDate](LayoutDate, v)
	return
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

func (t *CNTime) UnmarshalJSON(b []byte) (err error) {
	*t, err = unmarshalJSON[CNTime](LayoutN, b)
	return
}

func (t *CNTime) UnmarshalParam(param string) (err error) {
	*t, err = unmarshalParam[CNTime](LayoutN, param)
	return
}

func (t CNTime) Value() (driver.Value, error) {
	return value[CNTime](t)
}

func (t *CNTime) Scan(v any) (err error) {
	*t, err = scan[CNTime](LayoutN, v)
	return
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

func (t *CSTime) UnmarshalJSON(b []byte) (err error) {
	*t, err = unmarshalJSON[CSTime](LayoutS, b)
	return
}

func (t *CSTime) UnmarshalParam(param string) (err error) {
	*t, err = unmarshalParam[CSTime](LayoutS, param)
	return
}

func (t CSTime) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *CSTime) Scan(v any) (err error) {
	*t, err = scan[CSTime](LayoutS, v)
	return
}

func (CSTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return gormType(db)
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

func (t *CYM) UnmarshalJSON(b []byte) (err error) {
	*t, err = unmarshalJSON[CYM](LayoutYM, b)
	return
}

func (t *CYM) UnmarshalParam(param string) (err error) {
	*t, err = unmarshalParam[CYM](LayoutYM, param)
	return
}

func (t CYM) Value() (driver.Value, error) {
	return value[CYM](t)
}

func (t *CYM) Scan(v any) (err error) {
	*t, err = scan[CYM](LayoutYM, v)
	return
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
