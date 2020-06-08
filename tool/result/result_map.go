package result

import (
	"encoding/json"
	"log"
)

type ResultMap map[string]interface{}

func (c ResultMap) Add(key string, value interface{}) ResultMap {
	if c == nil {
		c = ResultMap{}
	}
	c[key] = value
	return c
}

func (c ResultMap) AddStruct(value interface{}) ResultMap {
	if c == nil {
		c = ResultMap{}
	}
	b, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil
	}
	//c[key] = value
	return c
}

// impl String()
func (c ResultMap) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Println("[ResultMap ERROR]:", err)
	}
	return string(b)
}

// 接口封装, 替换返回interface, 出现无法add操作问题
type Resultable interface {

	// 添加额外字段
	Add(key string, value interface{}) (rmp ResultMap)

	// 直接添加结构体
	AddStruct(value interface{}) (rmp ResultMap)
}
