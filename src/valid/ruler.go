package valid

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
