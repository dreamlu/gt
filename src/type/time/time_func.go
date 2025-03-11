package time

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func parse(layout, value string) (t time.Time, err error) {
	ll := len(layout)
	lv := len(value)
	if ll < lv {
		value = value[:ll]
	} else if ll > lv {
		value = fmt.Sprintf(`%s%s`, value, LayoutZero[lv:ll])
	}
	t, err = time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		value = fmt.Sprintf(`"%s"`, value)
		err = t.UnmarshalJSON([]byte(value))
		return
	}
	return
}

type ct interface {
	CTime | CDate | CNTime | CSTime | CYM
}

func marshalJSON[T ct](t T) ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t)), nil
}

func unmarshalParam[T ct](layout string, s string) (t T, err error) {
	return unmarshalJSON[T](layout, []byte(s))
}

func unmarshalJSON[T ct](layout string, b []byte) (t T, err error) {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return
	}
	ti, err := parse(layout, s)
	if err != nil {
		return t, err
	}
	t = T(ti)
	return
}

func value[T ct](t T) (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func scan[T ct](layout string, v any) (T, error) {
	var c T
	switch v.(type) {
	case []byte:
		value, err := parse(layout, string(v.([]byte)))
		if err != nil {
			return c, err
		}
		return T(value), nil
	case time.Time:
		value, ok := v.(time.Time)
		if ok {
			return T(value), nil
		}
	}
	return c, fmt.Errorf("can not scan [%v] format [%s]", v, layout)
}

func gormType(db *gorm.DB) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "timestamp" // timestamptz UTC
	default:
		return "datetime" // timestamp UTC
	}
}
