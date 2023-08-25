package mock

import (
	"fmt"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/json"
	"github.com/dreamlu/gt/src/type/time"
	"reflect"
	"strconv"
)

var (
	dataSize     = 21 // data size
	sliceSize    = 3  // slice size
	randomLength = 5  // char length
)

// SetDataSize set slice data size
func SetDataSize(size int) {
	dataSize = size
}

// SetSliceSize set data slice size
func SetSliceSize(size int) {
	sliceSize = size
}

// SetRandomLength set data random length
func SetRandomLength(len int) {
	randomLength = len
}

// Mock faker data
// support struct, []struct
// Support Only For:
// string,int,uint64,int64,float64,[]string,[]int,[]uint64,[]int64,[]float64
// CJSON,CTime,CDate
// other types will ignore
func Mock(data any) {
	typ := mr.TrueTypeof(data)
	if mr.IsStruct(typ) {
		mock(data)
	}
	if mr.IsSlice(typ) {
		ds := mr.ToSlice(data)
		for _, d := range ds {
			Mock(d)
		}
	}
}

func mock(data any) {
	var (
		typ      = mr.TrueTypeof(data)
		val      = mr.TrueValueOf(data)
		field    reflect.StructField
		fieldVal reflect.Value
	)
	for i := 0; i < typ.NumField(); i++ {
		field = typ.Field(i)
		fieldVal = val.Field(i)
		if field.Anonymous {
			mock(fieldVal.Interface())
			continue
		}
		var (
			value   any
			typName = field.Type.String() // field.Type.Name() can not get []string
		)
		value = mockValue(typName)
		if value == nil {
			continue
		}
		mr.Set(data, field.Name, value)
	}
}

func mockValue(typName string) (value any) {
	var (
		stringS = GetRand(Chinese, randomLength)
		numberS = GetRand(Number, randomLength)
		floatS  = fmt.Sprintf("%s.%s", numberS, GetRand(Number, 2))
	)

	switch typName {
	case "[]uint64":
		var tmp []uint64
		for k := 0; k < sliceSize; k++ {
			t, _ := strconv.ParseUint(numberS, 10, 64)
			tmp = append(tmp, t)
		}
		value = tmp
	case "uint64":
		value, _ = strconv.ParseUint(numberS, 10, 64)
	case "[]int64":
		var tmp []int64
		for k := 0; k < sliceSize; k++ {
			t, _ := strconv.ParseInt(numberS, 10, 64)
			tmp = append(tmp, t)
		}
		value = tmp
	case "int64":
		value, _ = strconv.ParseInt(numberS, 10, 64)
	case "[]int":
		var tmp []int
		for k := 0; k < sliceSize; k++ {
			t, _ := strconv.Atoi(numberS)
			tmp = append(tmp, t)
		}
		value = tmp
	case "int":
		value, _ = strconv.Atoi(numberS)
	case "[]float64":
		var tmp []float64
		for k := 0; k < sliceSize; k++ {
			t, _ := strconv.ParseFloat(floatS, 64)
			tmp = append(tmp, t)
		}
		value = tmp
	case "float64":
		value, _ = strconv.ParseFloat(floatS, 64)
	case "time.CDate":
		value = time.CDateNow()
	case "time.CTime":
		value = time.CTimeNow()
	case "json.CJSON":
		value = json.CJSON("{}")
	case "[]string":
		var tmp []string
		for k := 0; k < sliceSize; k++ {
			tmp = append(tmp, stringS)
		}
		value = tmp
	case "string":
		value = stringS
	default:
		value = nil
	}
	return
}
