package tmap

func Set[T TI](key string, value T) TMap[T] {
	return NewTMap[T]().Set(key, value)
}
