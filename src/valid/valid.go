package valid

import (
	"encoding/json"
	crudTag "github.com/dreamlu/gt/crud/dep/tag"
	mr "github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/tag"
	"github.com/dreamlu/gt/src/type/amap"
	"github.com/dreamlu/gt/src/type/cmap"
	"net/url"
	"reflect"
)

type Validator struct {
	// valid data
	data any
	// valid rule values
	rule ValidRule
}

// rule struct
type vRule struct {
	Vr Rule
	// required bool
}

type ValidSliceError struct {
	ValidError
	Line int
}

// valid type
type (
	ValidError map[string]error
	ValidRule  map[string]*vRule
)

func (v *ValidRule) String() string {
	s, _ := json.Marshal(v)
	return string(s)
}

func (v *ValidRule) Unmarshal(str string) {
	_ = json.Unmarshal([]byte(str), v)
}

func (v *ValidRule) parse(value any) {
	typ := reflect.TypeOf(value)
	for i := 0; i < typ.NumField(); i++ {
		// new rule
		rule := Rule{}
		// key
		rule.Key = tag.ParseJsonFieldTag(typ.Field(i)).Top()
		// rule

		field := typ.Field(i)
		rule.Valid = crudTag.ParseGtValidV(field)
		rule.Trans = crudTag.ParseGtTransV(field)
		if rule.Valid == "" {
			continue
		}
		// add rule
		(*v)[rule.Key] = &vRule{rule}
	}
}

var validBuffer = amap.NewAMap()

// Valid valid
func Valid(data any) ValidError {

	var typ reflect.Type
	typ, data = mr.TrueTypeofValue(data)

	if mr.IsSlice(typ) {
		if errs := validSlice(data, Valid); len(errs) > 0 {
			return errs[0].ValidError
		}
		return nil
	}

	return valid(data, typ)
}

func ValidSlice(data any) (sliceErrs []ValidSliceError) {
	return validSlice(data, Valid)
}

func validSlice(v any, vf func(data any) ValidError) (sliceErrs []ValidSliceError) {
	var (
		sls = mr.ToSlice(v)
	)
	for i, s := range sls {
		errs := vf(s)
		if len(errs) > 0 {
			sliceErrs = append(sliceErrs, ValidSliceError{
				ValidError: errs,
				Line:       i,
			})
		}
	}
	return
}

// ValidModel form/single json data
func ValidModel(data any, model any) ValidError {

	return valid(data, reflect.TypeOf(model))
}

// valid
func valid(value any, typ reflect.Type) ValidError {

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
				val any
				typ = reflect.TypeOf(d)
			)

			for i := 0; i < typ.NumField(); i++ {
				if crudTag.ParseJsonFieldTag(typ.Field(i)) == k {
					val = mr.Field(d, typ.Field(i).Name)
					break
				}
			}
			if err := r.Vr.Check(val); err != nil {
				errs[k] = err
			}
		}
	}
	return
}
