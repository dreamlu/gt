package valid

import (
	"encoding/json"
	"errors"
	"fmt"
	myReflect "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/te"
	"github.com/dreamlu/gt/tool/util/str"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validator struct {
	// 校验的数据对象
	data interface{}
	// 校检模型
	//model interface{}
	// 规则列表,key为字段名
	rule ValidRule
}

// 校检规则
type vRule struct {
	Vr DefaultRule
	// required bool
}

// 校验规则接口
// 支持自定义规则
type ValidateRuler interface {
	// 验证字段
	Check(data interface{}) error
}

// 内置规则结构
// 实现ValidateRuler接口
type DefaultRule struct {
	// 验证的字段名
	Key string
	// 翻译后的字段名
	// 默认 = Key
	Trans string
	// 规则
	Valid string
}

// valid error
type (
	ValidError map[string]error
	ValidRule  map[string]*vRule
)

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

var coMap = cmap.NewCMap()

// 创建校验器对象
// 针对结构体数据
// json
func Valid(data interface{}) ValidError {

	// 根据模型添加验证规则
	typ := reflect.TypeOf(data)
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

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
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
	//log.Println(key)
	vd := coMap.Get(key)
	if vd != "" {
		v.rule.Struct(vd)
		return v.Check()
	}

	v.rule = make(ValidRule)
	for i := 0; i < typ.NumField(); i++ {
		// 新建一个规则
		rule := DefaultRule{}
		// 字段名
		// 使用Name替代json Tag
		rule.Key = typ.Field(i).Tag.Get("json")
		// 规则
		gtTag := typ.Field(i).Tag.Get("gt")
		if gtTag == "" {
			continue
		}
		gtFields := strings.Split(gtTag, ";")
		for _, v := range gtFields {
			if strings.Contains(v, str.GtValid) {
				rule.Valid = string([]byte(v)[6:]) //strings.Trim(v, str.GtValid+":")
				//break
			}
			if strings.Contains(v, str.GtTrans) {
				rule.Trans = string([]byte(v)[6:]) //strings.TrimLeft(v, str.GtTrans+":")
				//break
			}
		}
		//rule.Valid = typ.Field(i).Tag.Get("valid")
		// 去除不存在验证字段
		if rule.Valid == "" {
			continue
		}
		// 字段翻译
		//rule.Trans = typ.Field(i).Tag.Get("Trans")

		// 绑定添加规则
		v.rule[rule.Key] = &vRule{rule}
	}

	// add coMap
	//log.Println(v.rule.String())
	coMap.Set(key, v.rule.String())

	// 进行校检
	return v.Check()
}

// 执行检查后返回信息
// Trans 翻译后的字段名
//func (v *Validator) CheckInfo() ValidError {
//	//if err := v.Check(); err != nil {
//	//	// 检查不通过，处理错误
//	//	// fmt.Println(err)
//	//	// return err
//	//	//for k := range v.rule {
//	//	//	if err[k] != nil {
//	//	//		return result.GetMapData(result.CodeValidator, err[k].Error())
//	//	//	}
//	//	//}
//	//}
//	return v.Check()
//}

// 执行检查
func (v *Validator) Check() (errs ValidError) {

	errs = make(ValidError)
	// 类型判断
	// d is value
	switch d := v.data.(type) {
	// request form
	// there is some duplicated: url.Values
	// maybe there is some ways to solve it!
	case url.Values:
		for k, r := range v.rule {
			// 数据
			data, _ := d[k]
			if data == nil {
				data = []string{""}
			}
			if err := r.Vr.Check(data[0]); err != nil { // 调用ValidateRuler接口的Check方法来检查
				errs[k] = err
			}
		}
	case cmap.CMap:
		for k, r := range v.rule {
			// 数据
			data, _ := d[k]
			if data == nil {
				data = []string{""}
			}
			if err := r.Vr.Check(data[0]); err != nil { // 调用ValidateRuler接口的Check方法来检查
				errs[k] = err
			}
		}
	default:
		for k, r := range v.rule {
			// 验证的字段值
			var val interface{}
			//val, err := myReflect.GetDataByFieldName(d, k)
			typ := reflect.TypeOf(d)
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			for i := 0; i < typ.NumField(); i++ {
				if typ.Field(i).Tag.Get("json") == k {
					//log.Println(reflect.ValueOf(d).Field(i).String())
					val, _ = myReflect.GetDataByFieldName(d, typ.Field(i).Name)
					break
				}
			}

			//log.Println(string(val))
			if err := r.Vr.Check(val); err != nil { // 调用ValidateRuler接口的Check方法来检查
				errs[k] = err
			}
		}
	}

	return errs
}

//  common rule
// 字段值转换成string进行验证
func (n *DefaultRule) Check(data interface{}) (err error) {
	// required 判断
	if !strings.Contains(n.Valid, "required") && data == "" {
		return
	}

	//  split rule
	//  先后规则顺序
	rules := strings.Split(n.Valid, ",")
	if n.Trans == "" {
		n.Trans = n.Key
	}
	for _, v := range rules {
		// 规则
		param := strings.Split(v, "=")

		switch param[0] {

		case "required":
			if err := nonzero(data); err != nil {
				return &te.TextError{Msg: fmt.Sprint(n.Trans, "为必填项")}
			}
		case "len":
			min := 0
			max := 0
			lg := length(data)
			// fix 中英文字符数量不统一
			// 范围判断
			if strings.Contains(param[1], "-") {
				args := strings.Split(param[1], "-")
				min, _ = strconv.Atoi(args[0])
				max, _ = strconv.Atoi(args[1])

			} else {
				max, _ = strconv.Atoi(param[1])
			}

			if lg < min || lg > max {
				return &te.TextError{Msg: fmt.Sprint(n.Trans, "长度在", min, "与", max, "之间")}
			}
		case "max":

			if err := max(data, param[1]); err != nil {
				return &te.TextError{Msg: fmt.Sprint(n.Trans, "最大值为", param[1])}
			}
		case "min":

			if err := min(data, param[1]); err != nil {
				return &te.TextError{Msg: fmt.Sprint(n.Trans, "最小值为", param[1])}
			}
		case "regex":
			if err := regex(data, param[1]); err != nil {
				return &te.TextError{Msg: fmt.Sprint(n.Trans, "正则规则为", param[1])}
			}
		case "phone":

			if b, _ := regexp.MatchString(`^1[2-9]\d{9}$`, data.(string)); !b {
				return &te.TextError{Msg: fmt.Sprintln("手机号格式非法")}
			}
		case "email":

			if b, _ := regexp.MatchString(`^([\w._]{2,10})@(\w+).([a-z]{2,4})$`, data.(string)); !b {
				return &te.TextError{Msg: fmt.Sprintln("邮箱格式非法")}
			}
		default:
			return errors.New(fmt.Sprintf("rule error: not support of rule=%s", n.Valid))
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

// string 类型, 切片数组判断长度
// 数值型应使用min和max判断大小
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

// 最小值
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

// 最大值
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

// 正则
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

//type myRuler struct {
//	// 验证的字段名
//	Key string
//	// 翻译后的字段名
//	// 默认 = Key
//	Trans string
//	// 规则
//	rule string
//}
//
//// 添加Check方法，实现ValidateRuler 接口
//func (m *myRuler) Check(data string) (Err error) {
//	// 判断data是否符合规则
//	return
//}
