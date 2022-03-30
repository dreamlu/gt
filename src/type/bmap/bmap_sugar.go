package bmap

func Set(key, value string) BMap {
	return NewBMap().Set(key, value)
}
