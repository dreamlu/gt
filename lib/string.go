package lib

func RemoveDuplicateString(languages []string) []string {
	var result []string
	temp := map[string]struct{}{}
	for _, item := range languages {
		// one key
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveStrings(strS []string, sep string) (res []string) {
	for _, v := range strS {
		if v != sep {
			res = append(res, v)
		}
	}
	return
}
