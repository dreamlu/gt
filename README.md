#### [个人]开发工具设计  
go-tool 是一个通用的api快速开发工具库  

##### 工具构成:  
| 路由    | orm  | 数据库 | 权限   |  配置   |  缓存  |
| ------ | ---- | ----  | ------ | ------ | ----- |
| go web | gorm | mysql | casbin(待完善)  | go-ini | redis |  

##### 原理：

1.封装  
2.反射  

##### demo(待完善)  
deercoder-gin  

##### 特点:
| 特点 | 
| ------ |
| 单/多张表的增删改查以及分页 |   
| 多张表连接操作 |  
| select `*`优化<br>反射替换*为具体字段名 |
| 优化自定义gorm日志<br>存储错误sql以及相关error |  
| 增加权限<br>用户-组(角色)-权限(菜单)(待优化) |
| 增加参数验证 |
| 增加mysql远程连接 |
| 增加多表key模糊搜索 |
| session(cookie/redis) |
| conf/app.conf 多开发模式支持 |
| 请求方式json/form data |
| [cache](./cache.go) 缓存实现 |
| [参数验证](tool/validator/validator_test.go) |  
| ...... |  

##### 使用  
go modules
##### API
- [API 使用](#api-examples)
    - [FORM请求](#form-request)
        - [Create](#create)
        - [Update](#update)
        - [Delete](#delete)
        - [GetBySearch](#getbysearch)
        - [GetByID](#getbyid)
        - [GetMoreBySearch](#getmorebysearch)
        - [GetDataBySQL](#getdatabysql)
        - [GetDataBySearchSQL](#getdatabysearchsql)
        - [DeleteBySQL](#deletebysql)
        - [UpdateBySQL](#updatebysql)
        - [CreateBySQL](#createbysql)
    - [JSON请求](#json-request)
    - [GetDevModeConfig](#getdevmodeconfig)
    - [缓存使用](#cachemanager)
    - [加解密](#aesende)
    - [标准日期](#time)
    - [JSON类型](#jsontype)
    - [字段验证](#validator)  
    - [日志支持](#customlog)
    

### API Examples  

#### Form Request

##### Create
```go
// dbcrud form data
var db = deercoder.DbCrud{
	Model: User{},		// model
	Table:"user",		// table name
}

// create user
func (c *User)Create(params map[string][]string) interface{} {

	params["createtime"] = append(params["createtime"], time.Now().Format("2006-01-02 15:04:05"))
	return db.Create(params)
}
```

##### Update
```go
// update user
func (c *User)Update(params map[string][]string) interface{} {

	return db.Update(params)
}
```

##### Delete
```go
// delete user, by id
func (c *User)Delete(id string) interface{} {

	return db.Delete(id)
}
```

##### GetBySearch
```go
// get user, limit and search
// clientPage 1, everyPage 10 default
func (c *User)GetBySearch(params map[string][]string) interface{} {
	var users []*User
	db.ModelData = &users
	return db.GetBySearch(params)
}
```

##### GetByID
```go
// get user, by id
func (c *User)GetByID(id string) interface{} {

	var user User	// not use *User
	db.ModelData = &user
	return db.GetByID(id)
}
```

##### GetMoreBySearch
```go
// get order, limit and search
// clientPage 1, everyPage 10 default
func (c *Order) GetMoreBySearch(params map[string][]string) interface{} {
	var or []*OrderD
	db = deercoder.DbCrud{
		InnerTables: []string{"order", "user"}, // inner join tables, 'order' must the first table
		LeftTables:  []string{"service"},       // left join tables
		Model:       OrderD{},                  // order model
		ModelData:   &or,                       // model value
	}
	return db.GetMoreBySearch(params)
}

```

##### GetDataBySQL
```go
// like UpdateBySQL
```

##### GetDataBySearchSQL
```go
// like UpdateBySQL
```

##### DeleteBySQL
```go
// like UpdateBySQL
```

##### UpdateBySQL
```go
var db = DbCrud{}
sql := "update `user` set name=? where id=?"
log.Println("[Info]:", db.UpdateBySQL(sql,"梦sql", 1))
```

##### CreateBySQL
```go
// like UpdateBySQL
```

#### Json Request
```go
// dbcrud json
// json request
var db_json = deercoder.DbCrudJ{
	Model: User{}, // model
	Table: "user", // table name
}

// get user, by id
func (c *User) GetByIDJ(id string) interface{} {

	var user User // not use *User
	db_json.ModelData = &user
	return db.GetByID(id)
}

// get user, limit and search
// clientPage 1, everyPage 10 default
func (c *User) GetBySearchJ(params map[string][]string) interface{} {
	var users []*User
	db_json.ModelData = &users
	return db.GetBySearch(params)
}

// delete user, by id
func (c *User) DeleteJ(id string) interface{} {

	return db_json.Delete(id)
}

// update user
func (c *User) UpdateJ(data *User) interface{} {

	return db_json.Update(data)
}

// create user
func (c *User) CreateJ(data *User) interface{} {

	// create time
	(*data).Createtime = deercoder.JsonTime(time.Now())

	return db_json.Create(data)
}

```

- 多模式配置文件  
> 配置方式: conf/app.conf 中 `devMode = dev` 对应conf/app-`dev`.conf  

#### GetDevModeConfig
```go
// devMode test
// app.conf devMode = dev
// test the app-dev.conf value
func TestDevMode(t *testing.T)  {
	log.Println("config read test: ", GetDevModeConfig("db.host"))
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


- 约定  
1.模型结构体json 内容与表字段保持一致  
2.返回格式参考[result](tool/result/result.go)    
3.多表关联命名, 模型中其他表字段命名: `他表名 + "_" + 他表字段名`  
n....  