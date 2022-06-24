package valid

import (
	"errors"
	"fmt"
	tErr "github.com/dreamlu/gt/src/type/errors"
)

// Ruler ruler
type Ruler func(rule string, data any) error

// RuleChain rule chain
// 设计模式--职责链模式,数组/map
type RuleChain map[string]Ruler

// 设计模式--单例模式[饿汉式]
var defaultRules RuleChain

func init() {
	defaultRules = map[string]Ruler{}
	addDefaultRuler()
}

// AddRuler add ruler
func AddRuler(key string, ruler Ruler) {
	defaultRules[key] = ruler
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
