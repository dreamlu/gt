package tmap

func Set[K comparable, V TI](key K, value V) TMap[K, V] {
	return NewTMap[K, V]().Set(key, value)
}
