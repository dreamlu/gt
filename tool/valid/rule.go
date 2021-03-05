package valid

import (
	"errors"
	"fmt"
	errors2 "github.com/dreamlu/gt/tool/type/errors"
	"regexp"
	"strconv"
	"strings"
)

// default rule type
//type defaultRule map[string]func(rule string, data interface{}) error

// add rule
func AddRule(key string, f func(rule string, data interface{}) error) {
	defaultRules[key] = f
}

func (n *Rule) rule(key, rule string, data interface{}) error {

	v, ok := defaultRules[key]
	if ok {
		err := v(rule, data)
		if err != nil {
			if errors.Is(err, errors2.TextErr) {
				return err
			}
			err = errors2.Text(fmt.Sprint(n.Trans, err.Error()))
		}
		return err
	}

	return errors.New(fmt.Sprintf("rule error: not support of rule=%s", n.Valid))
}

// default rules
var defaultRules = map[string]func(rule string, data interface{}) error{
	"required": func(rule string, data interface{}) error {
		if err := nonzero(data); err != nil {
			return errors.New("为必填项")
		}
		return nil
	},
	"len": func(rule string, data interface{}) error {
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
	},
	"max": func(rule string, data interface{}) error {
		if err := max(data, rule); err != nil {
			return errors.New(fmt.Sprint("最大值为", rule))
		}
		return nil
	},
	"min": func(rule string, data interface{}) error {
		if err := min(data, rule); err != nil {
			return errors.New(fmt.Sprint("最小值为", rule))
		}
		return nil
	},
	"regex": func(rule string, data interface{}) error {
		if err := regex(data, rule); err != nil {
			return errors.New(fmt.Sprint("正则规则为", rule))
		}
		return nil
	},
	"phone": func(rule string, data interface{}) error {
		if b, _ := regexp.MatchString(`^1[2-9]\d{9}$`, data.(string)); !b {
			return errors2.Text("手机号格式非法")
		}
		return nil
	},
	"email": func(rule string, data interface{}) error {
		if b, _ := regexp.MatchString(`^([\w._]{2,10})@(\w+).([a-z]{2,4})$`, data.(string)); !b {
			return errors2.Text("邮箱格式非法")
		}
		return nil
	},
}
