package param

type CMap map[string][]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v CMap) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (v CMap) Set(key, value string) {
	v[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v CMap) Add(key, value string) {
	v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v CMap) Del(key string) {
	delete(v, key)
}
