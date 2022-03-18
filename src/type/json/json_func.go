package json

import "encoding/json"

// CUnmarshal v to target
func CUnmarshal(v, t any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, t)
}
