package valid

import (
	"encoding/json"
	mr "github.com/dreamlu/gt/tool/reflect"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/util/cons"
	"github.com/dreamlu/gt/tool/util/tag"
	"log"
	"net/url"
	"reflect"
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

func (v *ValidRule) Unmarshal(str string) {
	err := json.Unmarshal([]byte(str), v)
	if err != nil {
		log.Println("valid string to struct err:", err)
	}
}

func (v *ValidRule) parse(value interface{}) {
	typ := reflect.TypeOf(value)
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
		(*v)[rule.Key] = &vRule{rule}
	}
}

var validBuffer = cmap.NewCMap()

// Valid valid
func Valid(data interface{}) ValidError {

	var typ reflect.Type
	typ, data = mr.TrueTypeofValue(data)

	if typ.Kind() == reflect.Slice {
		if errs := validSlice(data, Valid); len(errs) > 0 {
			return errs
		}
		return nil
	}

	return valid(data, typ)
}

func validSlice(v interface{}, vf func(data interface{}) ValidError) ValidError {
	var (
		sls  = mr.ToSlice(v)
		errs ValidError
	)
	for _, s := range sls {
		errs = vf(s)
		if len(errs) > 0 {
			return errs
		}
	}
	return nil
}

// ValidModel form/single json data
func ValidModel(data interface{}, model interface{}) ValidError {

	return valid(data, reflect.TypeOf(model))
}

// valid
func valid(value interface{}, typ reflect.Type) ValidError {

	var (
		v   = &Validator{data: value}
		key = mr.Path(typ, "valid")
		vd  = validBuffer.Get(key)
	)

	if vd != "" {
		v.rule.Unmarshal(vd)
		return v.Check()
	}

	v.rule = make(ValidRule)
	v.rule.parse(value)
	if len(v.rule) == 0 {
		return nil
	}
	validBuffer.Set(key, v.rule.String())

	return v.Check()
}

// Check rule
func (v *Validator) Check() (errs ValidError) {

	errs = make(ValidError)
	// d is value
	switch d := v.data.(type) {
	// request form
	// there is some duplicated: url.Values
	// maybe there is some ways to solve it!
	case url.Values, cmap.CMap:
		if vd, ok := v.data.(url.Values); ok {
			d = cmap.CMap(vd)
		}
		for k, r := range v.rule {
			data, _ := d.(cmap.CMap)[k]
			if data == nil {
				data = []string{""}
			}
			if err := r.Vr.Check(data[0]); err != nil {
				errs[k] = err
			}
		}
	default:
		for k, r := range v.rule {
			var (
				val interface{}
				typ = reflect.TypeOf(d)
			)

			for i := 0; i < typ.NumField(); i++ {
				if tag.GetFieldTag(typ.Field(i)) == k {
					val, _ = mr.ValueOfName(d, typ.Field(i).Name)
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

// Check rule common rule Check
func (n *Rule) Check(data interface{}) (err error) {
	// required
	if !strings.Contains(n.Valid, RuleRequired) && data == "" {
		return
	}

	//  split rule
	rules := strings.Split(n.Valid, cons.GtComma)
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
