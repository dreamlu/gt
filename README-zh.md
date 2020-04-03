#### Erwin Schrödinger's Cat  
gt使用手册 (v1.7.0)  

api快速开发工具库,模型生成  

##### demo: [deercoder-gin](https://github.com/dreamlu/deercoder-gin)  

##### API
- [API 使用](#api-examples)
    - [模型定义](#model)
    - [Crud](#Crud-request)
    - [批量创建](#createmore)
    - [配置文件模式](#getdevmode)
    - [缓存使用](#cachemanager)
    - [加解密](#aesende)
    - [标准日期](#time)
    - [JSON类型](#jsontype)
    - [字段验证](#validator)  
    - [日志支持](#customlog)
    - [snowflake ID](#snowflakeid)
- [扩展 使用](#extend-examples)
    - [crud原生SQL](#crud-selectsql)
    - [更新其他字段](#crud-update)
    - [事务](#transcation)
    

### API Examples  

#### Model
```go
// user model
type User struct {
	ID         uint64     `json:"id"`
	Name       string     `json:"name"`
	Createtime time.CTime `json:"createtime"`
}

// service model
type Service struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// order model
type Order struct {
	ID         uint64     `json:"id"`
	UserID     int64      `json:"user_id"`    // user id
	ServiceID  int64      `json:"service_id"` // service table id
	Createtime time.CTime `json:"createtime"` // createtime
}

// order detail
type OrderD struct {
	Order
	UserName    string `json:"user_name"`    // user table column name
	ServiceName string `json:"service_name"` // service table column `name`
}

    // select more
    // 多表查询
	// get more search
	var params = make(cmap.CMap)
	params.Add("user_id", "1")
	//params.Add("key", "梦") // key word
	params.Add("clientPage", "1") // 第一页
	params.Add("everyPage", "2") // 每页2条
	var or []*OrderD
	crud := NewCrud(
		InnerTable([]string{"order", "user", "order", "service"}),
		//LeftTable([]string{"order", "service"}),
		Model(OrderD{}),
		Data(&or),
		SubWhereSQL("1 = 1", "2 = 2", ""),
	)
	err := crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
	t.Log("\n[User Info]:", or[0])
    
// output:
    TestGetMoreDataBySearch: db_test.go:243: 
        [User Info]: &{{1 1 1 "2019-01-28 15:07:06"} 梦sql 服务名称}
--- PASS: TestGetMoreDataBySearch (0.00s)
PASS
```

#### Crud Request
```go
// GET Request
//func ToCMap(u *gin.Context) cmap.CMap {
//	err := u.Request.ParseForm()
//	if err != nil {
//		gt.Logger().Error(err.Error())
//		return nil
//	}
//	values := cmap.CMap(u.Request.Form)
//	xss.XssMap(values)
//	return values
//}

var crud = gt.NewCrud(
	gt.Model(Client{}),
)

// get data, by id
func (c *Client) GetByID(id string) (*Client, error) {

	var data Client // not use *Client
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Client) GetBySearch(params cmap.CMap) (datas []*Client, pager result.Pager, err error) {
	//var datas []*Client
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *Client) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *Client) Update(data *Client) (*Client, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// create data
func (c *Client) Create(data *Client) (*Client, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
```

#### CreateMore
```go
// 批量创建
func TestCreateMoreDataJ(t *testing.T) {
    type UserPar struct {
		Name       string     `json:"name"`
		Createtime time.CTime `json:"createtime"`
	}
	type User struct {
		ID uint64 `json:"id"`
		UserPar
	}
	
	var up = []UserPar{
		{Name: "测试1", Createtime: time.CTime(time2.Now())},
		{Name: "测试2"},
	}
	crud := NewCrud(
		//Table("user"),
		Model(UserPar{}),
		Data(up),
		//SubSQL("(asdf) as a","(asdfa) as b"),
	)

	err := crud.CreateMoreData()
	t.Log(err)
}

```

- 多模式配置文件  
> 配置方式: conf/app.yaml 中 `devMode = dev` 对应conf/app-`dev`.yaml  

#### GetDevMode
```go
    type dba struct {
    	User        string
    	Password    string
    	Host        string
    	Name        string
    	MaxIdleConn int
    	MaxOpenConn int
    	// db log mode
    	Log bool
    }
    config := gt.Configger()
    dbS := &dba{
        user:     config.GetString("app.db.user"),
        password: config.GetString("app.db.password"),
        host:     config.GetString("app.db.host"),
        name:     config.GetString("app.db.name"),
    }
    // or
    dbS := &dba{}
    gt.Configger().GetStruct("app.db", dbS)
```

#### CacheManager
```go
    ce = gt.NewCache()
    data := CacheModel{
		Time: 50 * CacheMinute,
		Data: user,
	}

	// key can use user.ID,user.Name,user
	// because it can be interface
	// set
	err := ce.Set(user, data)
	t.Log("set err: ", err)

	// get
	reply, _ := ce.Get(user)
	t.Log("user data :", reply.Data)
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
type CJSON []byte
```  

### Validator  
```go
func TestValidator(t *testing.T) {

	type Test struct {
		ID   int64  `json:"id" valid:"required,min=0,max=5"`
		Name string `json:"name" valid:"required,len=2-5" trans:"用户名"`
	}

	// form request data
	var maps = make(map[string][]string)
	maps["name"] = append(maps["name"], "梦1")
	info := Valid(maps, Test{})
	log.Println(info)

	// json request data
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

### Extend Examples  

#### Crud Selectsql
```go
    sql := "update `user` set name=? where id=?"
    t.Log("[Info]:", crud.Select(sql, "梦sql", 1).Exec())
    var user []*User
	crud.Params(
		Data(&user),
		ClientPage(1),
		EveryPage(2),
	).
		Select("select *from user").
		Select("where id > 0")
	if true {
		crud.Select("and 1=1")
	}
	crud.Search() // 查询 + 分页
    t.Log(crud.Pager()) // Search()分页数据
	crud.Single() // 直接查询
```

#### Crud Update
```go
    type UserPar struct {
		Name string `json:"name"`
	}
	crud := crud.Params(
		//Table("user"),
		Model(User{}),
		Data(&UserPar{
			//ID:   1,
			Name: "梦S",
		}),
	)
	t.Log(crud.Update().RowsAffected())
	t.Log(crud.Select("`name` = ?", "梦").Update().RowsAffected())
	t.Log(crud.Error())
```

#### Transcation
```go
    cd := crud.Begin()
	cd.Params(
		Table("user"),
		Data(&User{
			ID:   11234,
			Name: "梦S",
		}),
	).Create()
	if cd.Error() != nil {
		cd.Rollback()
	}
	cd.Params(
		Data(&User{
			Name: "梦SSS2",
		})).Create()
	if cd.Error() != nil {
		cd.Rollback()
	}
	// add select sql test
	var u []User
	cd.Params(Data(&u)).Select("select * from `user`").Select("where 1=1").Single()
	cd.Params(Data(&u)).Select("select * from `user`").Select("where 1=1").Single()
	//cd.DB().Raw("select * from `user`").Scan(&u)

	cd.Commit()
	if cd.Error() != nil {
		cd.Rollback()
	}
```

- 约定  
1.模型结构体json 内容与表字段保持一致  
2.返回格式参考[result](tool/result/result.go)    
3.多表关联命名, 模型中其他表字段命名: `他表名 + "_" + 他表字段名`  
n....  