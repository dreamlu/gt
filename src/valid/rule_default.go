package valid

import (
	"errors"
	"fmt"
	errors2 "github.com/dreamlu/gt/src/type/errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	RuleRequired = "required"
	RuleLen      = "len"
	RuleMax      = "max"
	RuleMin      = "min"
	RuleRegex    = "regex"
	RulePhone    = "phone"
	RuleEmail    = "email"
)

func addDefaultRuler() {
	AddRuler(RuleRequired, requiredRuler)
	AddRuler(RuleLen, lenRuler)
	AddRuler(RuleMax, maxRuler)
	AddRuler(RuleMin, minRuler)
	AddRuler(RuleRegex, regexRuler)
	AddRuler(RulePhone, phoneRuler)
	AddRuler(RuleEmail, emailRuler)
}

func requiredRuler(rule string, data any) error {
	if err := nonzero(data); err != nil {
		return errors.New("为必填项")
	}
	return nil
}

func lenRuler(rule string, data any) error {
	min := 0
	max := 0
	lg := length(data)
	if strings.Contains(rule, "-") {
		args := strings.Split(rule, "-")
		min, _ = strconv.Atoi(args[0])
		max, _ = strconv.Atoi(args[1])

	} else {
		max, _ = strconv.Atoi(rule)
	}

	if lg < min || lg > max {
		return errors.New(fmt.Sprint("长度在", min, "与", max, "之间"))
	}
	return nil
}

func maxRuler(rule string, data any) error {
	if err := max(data, rule); err != nil {
		return errors.New(fmt.Sprint("最大值为", rule))
	}
	return nil
}

func minRuler(rule string, data any) error {
	if err := min(data, rule); err != nil {
		return errors.New(fmt.Sprint("最小值为", rule))
	}
	return nil
}

func regexRuler(rule string, data any) error {
	if err := regex(data, rule); err != nil {
		return errors.New(fmt.Sprint("正则规则为", rule))
	}
	return nil
}

func phoneRuler(rule string, data any) error {
	if b, _ := regexp.MatchString(`^1[2-9]\d{9}$`, data.(string)); !b {
		return errors2.Text("手机号格式非法")
	}
	return nil
}

func emailRuler(rule string, data any) error {
	if b, _ := regexp.MatchString(`^([\w._]{2,10})@(\w+).([a-z]{2,4})$`, data.(string)); !b {
		return errors2.Text("邮箱格式非法")
	}
	return nil
}

// default ruler func

func nonzero(v any) error {
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
func length(v any) int {
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
	default:
		return 0
	}
}

// min
func min(v any, param string) error {
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
func max(v any, param string) error {
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
func regex(v any, param string) error {
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

// asInt return the parameter as a int64
// or panics if it can't convert
func asInt(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// asUint return the parameter as a uint64
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
