package validator

import (
	"errors"
	"fmt"
	myReflect "github.com/dreamlu/go-tool/tool/reflect"
	"github.com/dreamlu/go-tool/tool/result"
	"github.com/dreamlu/go-tool/tool/type/te"
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
	model interface{}
	// 规则列表,key为字段名
	rule map[string]*vRule
}

// 校检规则
type vRule struct {
	vr ValidateRuler
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
	key string
	// 翻译后的字段名
	// 默认 = key
	trans string
	// 规则
	valid string
}

// 创建校验器对象
// 针对form表单/json数据两种
func Valid(data, model interface{}) result.MapData {

	v := &Validator{
		data:  data,
		model: model,
	}
	v.rule = make(map[string]*vRule)
	// 根据模型添加验证规则
	typ := reflect.TypeOf(model)
	for i := 0; i < typ.NumField(); i++ {
		// 新建一个规则
		rule := &DefaultRule{}
		// 字段名
		// 使用Name替代json Tag
		rule.key = typ.Field(i).Tag.Get("json")
		// 规则
		rule.valid = typ.Field(i).Tag.Get("valid")
		// 去除不存在验证字段
		if rule.valid == "" {
			continue
		}
		// 字段翻译
		rule.trans = typ.Field(i).Tag.Get("trans")

		// 绑定添加规则
		v.rule[rule.key] = &vRule{rule}
	}
	// 进行校检
	return v.CheckInfo()
}

// 执行检查后返回信息
// trans 翻译后的字段名
func (v *Validator) CheckInfo() result.MapData {
	if err := v.Check(); err != nil {
		// 检查不通过，处理错误
		// fmt.Println(err)
		// return err
		for k := range v.rule {
			if err[k] != nil {
				return result.GetMapData(result.CodeValidator, err[k].Error())
			}
		}
	}
	return result.MapValSuccess
}

// 执行检查
func (v *Validator) Check() (errs map[string]error) {

	errs = make(map[string]error)
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
			if err := r.vr.Check(data[0]); err != nil { // 调用ValidateRuler接口的Check方法来检查
				errs[k] = err
			}
		}
	case map[string][]string:
		for k, r := range v.rule {
			// 数据
			data, _ := d[k]
			if data == nil {
				data = []string{""}
			}
			if err := r.vr.Check(data[0]); err != nil { // 调用ValidateRuler接口的Check方法来检查
				errs[k] = err
			}
		}
	default:
		for k, r := range v.rule {
			// 验证的字段值
			var val interface{}
			//val, err := myReflect.GetDataByFieldName(d, k)
			typ := reflect.TypeOf(d)
			for i := 0; i < typ.NumField(); i++ {
				if typ.Field(i).Tag.Get("json") == k {
					//log.Println(reflect.ValueOf(d).Field(i).String())
					val, _ = myReflect.GetDataByFieldName(d, typ.Field(i).Name)
					break
				}
			}

			//log.Println(string(val))
			if err := r.vr.Check(val); err != nil { // 调用ValidateRuler接口的Check方法来检查
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
	if !strings.Contains(n.valid, "required") && data == "" {
		return
	}

	//  split rule
	//  先后规则顺序
	rules := strings.Split(n.valid, ",")
	if n.trans == "" {
		n.trans = n.key
	}
	for _, v := range rules {
		// 规则
		param := strings.Split(v, "=")

		switch param[0] {

		case "required":
			if err := nonzero(data); err != nil {
				return &te.TextError{Msg: fmt.Sprintln(n.trans, "为必填项")}
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
				return &te.TextError{Msg: fmt.Sprintln(n.trans, "长度在", min, "与", max, "之间")}
			}
		case "max":

			if err := max(data, param[1]); err != nil {
				return &te.TextError{Msg: fmt.Sprintln(n.trans, "最大值为", param[1])}
			}
		case "min":

			if err := min(data, param[1]); err != nil {
				return &te.TextError{Msg: fmt.Sprintln(n.trans, "最小值为", param[1])}
			}
		case "regex":
			if err := regex(data, param[1]); err != nil {
				return &te.TextError{Msg: fmt.Sprintln(n.trans, "正则规则为", param[1])}
			}
		case "phone":

			if b, _ := regexp.MatchString(`^1[2-9]\d{9}$`, data.(string)); !b {
				return &te.TextError{Msg: fmt.Sprintln("手机号格式非法")}
			}
		case "email":

			if b, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w+).([a-z]{2,4})$`, data.(string)); !b {
				return &te.TextError{Msg: fmt.Sprintln("邮箱格式非法")}
			}
		default:
			return errors.New(fmt.Sprintf("rule error: not support of rule=%s", n.valid))
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
		return errors.New(result.MsgValError)
	}

	if !valid {
		return errors.New(result.MsgValError)
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
			return errors.New(result.MsgValError)
		}
		invalid = int64(len(st.String())) < p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = int64(st.Len()) < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = st.Int() < p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = st.Uint() < p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = st.Float() < p
	default:
		return errors.New(result.MsgValError)
	}
	if invalid {
		return errors.New(result.MsgValError)
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
			return errors.New(result.MsgValError)
		}
		invalid = int64(len(st.String())) > p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = int64(st.Len()) > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = st.Int() > p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = st.Uint() > p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return errors.New(result.MsgValError)
		}
		invalid = st.Float() > p
	default:
		return errors.New(result.MsgValError)
	}
	if invalid {
		return errors.New(result.MsgValError)
	}
	return nil
}

// 正则
func regex(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		sptr, ok := v.(*string)
		if !ok {
			return errors.New(result.MsgValError)
		}
		if sptr == nil {
			return nil
		}
		s = *sptr
	}

	re, err := regexp.Compile(param)
	if err != nil {
		return errors.New(result.MsgValError)
	}

	if !re.MatchString(s) {
		return errors.New(result.MsgValError)
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
//	key string
//	// 翻译后的字段名
//	// 默认 = key
//	trans string
//	// 规则
//	rule string
//}
//
//// 添加Check方法，实现ValidateRuler 接口
//func (m *myRuler) Check(data string) (Err error) {
//	// 判断data是否符合规则
//	return
//}
