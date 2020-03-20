package result

type ResultMap map[string]interface{}

func (c ResultMap) Add(key string, value interface{}) ResultMap {
	if c == nil {
		c = ResultMap{}
	}
	c[key] = value
	return c
}
