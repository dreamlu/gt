package result

import "encoding/json"

type ResultMap map[string]interface{}

func (c ResultMap) Add(key string, value interface{}) ResultMap {
	if c == nil {
		c = ResultMap{}
	}
	c[key] = value
	return c
}

func (c ResultMap) AddStruct(value interface{}) ResultMap {
	if c == nil {
		c = ResultMap{}
	}
	b, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil
	}
	//c[key] = value
	return c
}
