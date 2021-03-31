package valid

import (
	"encoding/json"
	"errors"
	"fmt"
	myReflect "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/tag"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/cons"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validator struct {
	// valid data
	data interface{}
	// valid rule values
	rule ValidRule
}

// rule struct
type vRule struct {
	Vr Rule
	// required bool
}

// valid type
type (
	ValidError map[string]error
	ValidRule  map[string]*vRule
)

type Rule struct {
	// key
	Key string
	// 翻译后的字段名
	// default Key
	Trans string
	// valid
	Valid string
}

func (v *ValidRule) String() string {
	s, err := json.Marshal(v)
	if err != nil {
		log.Println("valid struct:[", err, "]:error")
		return ""
	}
	return string(s)
}

func (v *ValidRule) Struct(str string) {
	err := json.Unmarshal([]byte(str), v)
	if err != nil {
		log.Println("valid string to struct err:", err)
	}
}

var validBuffer = cmap.NewCMap()

// Valid
func Valid(data interface{}) ValidError {

	typ := reflect.TypeOf(data)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		data = reflect.ValueOf(data).Elem().Interface()
	}

	if typ.Kind() == reflect.Slice {
		ss := myReflect.ToSlice(data)
		for _, v := range ss {
			errs := Valid(v)
			if len(errs) > 0 {
				return errs
			}
		}
		return nil
	}

	return valid(data, typ)
}

// form/single json data
func ValidModel(data interface{}, model interface{}) ValidError {

	return valid(data, reflect.TypeOf(model))
}

// valid
func valid(data interface{}, typ reflect.Type) ValidError {

	var (
		v = &Validator{data: data}
	)
	key := typ.PkgPath() + "/valid/" + typ.Name()
	vd := validBuffer.Get(key)
	if vd != "" {
		v.rule.Struct(vd)
		return v.Check()
	}

	v.rule = make(ValidRule)
	for i := 0; i < typ.NumField(); i++ {
		// new rule
		rule := Rule{}
		// key
		rule.Key = tag.GetFieldTag(typ.Field(i))
		// rule
		gtTag := typ.Field(i).Tag.Get(cons.GT)
		if gtTag == "" {
			continue
		}
		gtFields := strings.Split(gtTag, ";")
		for _, v := range gtFields {
			if strings.Contains(v, cons.GtValid) {
				rule.Valid = string([]byte(v)[6:])
			}
			if strings.Contains(v, cons.GtTrans) {
				rule.Trans = string([]byte(v)[6:])
			}
		}
		if rule.Valid == "" {
			continue
		}
		// add rule
		v.rule[rule.Key] = &vRule{rule}
	}

	validBuffer.Set(key, v.rule.String())

	return v.Check()
}

// Check
func (v *Validator) Check() (errs ValidError) {

	errs = make(ValidError)
	// d is value
	switch d := v.data.(type) {
	// request form
	// there is some duplicated: url.Values
	// maybe there is some ways to solve it!
	case url.Values:
		for k, r := range v.rule {
			data, _ := d[k]
			if data == nil {
				data = []string{""}
			}
			if err := r.Vr.Check(data[0]); err != nil {
				errs[k] = err
			}
		}
	case cmap.CMap:
		for k, r := range v.rule {
			data, _ := d[k]
			if data == nil {
				data = []string{""}
			}
			if err := r.Vr.Check(data[0]); err != nil {
				errs[k] = err
			}
		}
	default:
		for k, r := range v.rule {
			var val interface{}
			typ := reflect.TypeOf(d)
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			for i := 0; i < typ.NumField(); i++ {
				if tag.GetFieldTag(typ.Field(i)) == k {
					val, _ = myReflect.GetDataByFieldName(d, typ.Field(i).Name)
					break
				}
			}

			if err := r.Vr.Check(val); err != nil {
				errs[k] = err
			}
		}
	}

	return errs
}

//  rule common rule Check
func (n *Rule) Check(data interface{}) (err error) {
	// required
	if !strings.Contains(n.Valid, "required") && data == "" {
		return
	}

	//  split rule
	rules := strings.Split(n.Valid, ",")
	if n.Trans == "" {
		n.Trans = n.Key
	}
	for _, v := range rules {
		param := strings.Split(v, "=")
		rule := ""
		if length(param) > 1 {
			rule = param[1]
		}
		if err = n.rule(param[0], rule, data); err != nil {
			return err
		}
	}
	return nil
}

func nonzero(v interface{}) error {
	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		valid = len(st.String()) != 0
	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		valid = st.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() != 0
	case reflect.Float32, reflect.Float64:
		valid = st.Float() != 0
	case reflect.Bool:
		valid = st.Bool()
	case reflect.Invalid:
		valid = false // always invalid
	case reflect.Struct:
		valid = true // always valid since only nil pointers are empty
	default:
		return fmt.Errorf("%w", errors.New("[nonzero val error]"))
	}

	if !valid {
		return fmt.Errorf("%w", errors.New("[nonzero val error]"))
	}
	return nil
}

// length
func length(v interface{}) int {
	st := reflect.ValueOf(v)
	if st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return 0
		}
		st = st.Elem()
	}
	switch st.Kind() {
	case reflect.String:
		return len([]rune(st.String()))
	case reflect.Slice, reflect.Map, reflect.Array:
		return st.Len()
	//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//	log.Println(st)
	//	return len(st)
	//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	//	return st.Len()
	//case reflect.Float32, reflect.Float64:
	//	st.Len()
	default:
		return 0
	}
	//return 0
}

// min
func min(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	invalid := false
	if st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return nil
		}
		st = st.Elem()
	}
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[min val error]"))
		}
		invalid = int64(len(st.String())) < p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[min val error]"))
		}
		invalid = int64(st.Len()) < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[min val error]"))
		}
		invalid = st.Int() < p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[min val error]"))
		}
		invalid = st.Uint() < p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[min val error]"))
		}
		invalid = st.Float() < p
	default:
		return fmt.Errorf("%w", errors.New("[min val error]"))
	}
	if invalid {
		return fmt.Errorf("%w", errors.New("[min val error]"))
	}
	return nil
}

// max
func max(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	var invalid bool
	if st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return nil
		}
		st = st.Elem()
	}
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[max val error]"))
		}
		invalid = int64(len(st.String())) > p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[max val error]"))
		}
		invalid = int64(st.Len()) > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[max val error]"))
		}
		invalid = st.Int() > p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[max val error]"))
		}
		invalid = st.Uint() > p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return fmt.Errorf("%w", errors.New("[max val error]"))
		}
		invalid = st.Float() > p
	default:
		return fmt.Errorf("%w", errors.New("[max val error]"))
	}
	if invalid {
		return fmt.Errorf("%w", errors.New("[max val error]"))
	}
	return nil
}

// regex
func regex(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		sptr, ok := v.(*string)
		if !ok {
			return fmt.Errorf("%w", errors.New("[regex val error]"))
		}
		if sptr == nil {
			return nil
		}
		s = *sptr
	}

	re, err := regexp.Compile(param)
	if err != nil {
		return fmt.Errorf("%w", errors.New("[regex val error]"))
	}

	if !re.MatchString(s) {
		return fmt.Errorf("%w", errors.New("[regex val error]"))
	}
	return nil
}

// asInt retuns the parameter as a int64
// or panics if it can't convert
func asInt(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// asUint retuns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) (uint64, error) {
	i, err := strconv.ParseUint(param, 0, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// asFloat retuns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) (float64, error) {
	i, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0.0, err
	}
	return i, nil
}
