package valid

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/crud/dep/cons"
	tErr "github.com/dreamlu/gt/src/type/errors"
	"strings"
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

// Check rule common rule Check
func (n *Rule) Check(data any) (err error) {
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

func (n *Rule) rule(key, rule string, data any) error {

	v, ok := defaultRules[key]
	if ok {
		err := v(rule, data)
		if err != nil {
			if errors.Is(err, tErr.TextErr) {
				return err
			}
			err = tErr.Text(fmt.Sprint(n.Trans, err.Error()))
		}
		return err
	}

	return errors.New(fmt.Sprintf("rule error: not support of rule=%s", key))
}
