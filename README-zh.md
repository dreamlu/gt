#### Erwin Schrödinger's Cat  
gt(gt) (old)  

web快速开发工具库,模型生成  

##### demo  
[deercoder-gin](https://github.com/dreamlu/deercoder-gin)  

##### API
- [API 使用](#api-examples)
    - [SQL动态请求](#SQL-request)
    - [批量创建](#createmore)
    - [配置文件模式](#getdevmode)
    - [缓存使用](#cachemanager)
    - [加解密](#aesende)
    - [标准日期](#time)
    - [JSON类型](#jsontype)
    - [字段验证](#validator)  
    - [日志支持](#customlog)
    - [snowflake ID](#snowflakeid)
    

### API Examples  

#### SQL Request

###### 如何使用
参考[测试](db_test.go)

#### CreateMore
```go
// 批量创建
func TestCreateMoreDataJ(t *testing.T) {

	var user = []User{
		{Name: "测试1", Createtime: time.CTime(time2.Now())},
		{Name: "测试2"},
	}

	crud.Params = &Params{
		Table: "user",
		Model: User{},
	}

	err := crud.CreateMoreData(user)
	log.Println(err)
}

```

- 多模式配置文件  
> 配置方式: conf/app.yaml 中 `devMode = dev` 对应conf/app-`dev`.yaml  

#### GetDevMode
```go
	config := &Config{}
	config.NewConfig()

	dbS := &dba{
		user:     config.GetString("app.db.user"),
		password: config.GetString("app.db.password"),
		host:     config.GetString("app.db.host"),
		name:     config.GetString("app.db.name"),
	}
```

#### CacheManager
```go
var cache CacheManager = new(RedisManager)

func init()  {
	// init redis
	//_ = r.NewCache()
	// init cache
	_ = cache.NewCache()
}

// 具体操作
// cache manager
type CacheManager interface {
	// init cache
	NewCache() error
	// operate method
	// set value
	// if time != 0 set it
	Set(key interface{}, value CacheModel) error
	// get value
	Get(key interface{}) (CacheModel, error)
	// delete value
	Delete(key interface{}) error
	// more del
	// key will become *key*
	DeleteMore(key interface{}) error
	// check value
	// flush the time
	Check(key interface{}) error
}
```

### AesEnDe  
```go
log.Println("[加密测试]:", AesEn("123456"))
log.Println("[解密测试]:", AesDe("lIEbR7cEp2U10gtM0j8dCg=="))
```

### Time
```go
// 时间格式化2006-01-02 15:04:05
type CTime time.Time
// 时间格式化2006-01-02
type CDate time.Time 
```  

### JSONType
```go
// 返回json类型
type CJSON time.Time
```  

### Validator  
```go
func TestValidator(t *testing.T) {

	type Test struct {
		ID   int64  `json:"id" valid:"required,min=0,max=5"`
		Name string `json:"name" valid:"required,len=2-5" trans:"用户名"`
	}

	// form data
	var maps = make(map[string][]string)
	maps["name"] = append(maps["name"], "梦1")
	info := Valid(maps, Test{})
	log.Println(info)

	// json data
	var test = Test{
		ID:   6,
		Name: "梦1",
	}
	log.Println(Valid(test, Test{}))

}
```

### CustomLog
```go
func TestNewFileLog(t *testing.T) {

	myLog.Info("项目路径", projectPath)
	for {
		time.Sleep(1 * time.Second)
		myLog.Error("测试")
	}
}
```

### SnowflakeId
```go
func TestId(t *testing.T) {
	id, err := NewID(1)
	if err != nil {
		log.Print(err)
		return
	}
	t.Log(id.String())
}
```


- 约定  
1.模型结构体json 内容与表字段保持一致  
2.返回格式参考[result](tool/result/result.go)    
3.多表关联命名, 模型中其他表字段命名: `他表名 + "_" + 他表字段名`  
n....  