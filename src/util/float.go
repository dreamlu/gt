package util

import (
	"fmt"
	"strconv"
)

func Float64(value float64, size int) float64 {
	var (
		format = `%.` + strconv.Itoa(size) + `f`
	)
	value, _ = strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return value
}
