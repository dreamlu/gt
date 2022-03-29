package cmap

func Set(key, value string) CMap {
	return NewCMap().Set(key, value)
}

func Add(key, value string) CMap {
	return NewCMap().Add(key, value)
}
