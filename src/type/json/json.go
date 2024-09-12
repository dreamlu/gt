package json

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type CJSON []byte

func (j CJSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

func (j *CJSON) Scan(v any) error {
	if v == nil {
		*j = nil
		return nil
	}
	var s []byte
	switch v.(type) {
	case string:
		s = []byte(v.(string))
	case []byte:
		s = v.([]byte)
	default:
		return fmt.Errorf("[not json data error]:%s", v)
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (CJSON) GormDataType() string {
	return "json"
}

func (j CJSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		// use "" replace null
		return []byte("\"\""), nil
	}
	return j, nil
}

func (j *CJSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("CJSON nil error")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j CJSON) Equals(j1 CJSON) bool {
	return bytes.Equal(j, j1)
}

func (j CJSON) String() string {
	return string(j)
}

// Unmarshal support Struct/Array
func (j CJSON) Unmarshal(v any) error {
	err := json.Unmarshal(j, v)
	if err != nil {
		return err
	}
	return nil
}
