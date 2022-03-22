package lib

func RemoveDuplicate[T comparable](s []T) []T {
	var result []T
	temp := map[T]struct{}{}
	for _, item := range s {
		// one key
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Remove[T comparable](s []T, sep ...T) (res []T) {
	for _, v := range s {
		for _, se := range sep {
			if v != se {
				res = append(res, v)
			}
		}
	}
	return
}

func Contains[T comparable](s []T, sep ...T) bool {
	for _, v := range s {
		for _, se := range sep {
			if v == se {
				return true
			}
		}
	}
	return false
}
