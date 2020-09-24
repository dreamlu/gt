#### Erwin Schrödinger's Cat  
gt使用手册 (v1.10.0+)  

api快速开发业务框架,模型生成  
通用增删改查，支持多表连接  

##### demo:  
[deercoder-gin](https://github.com/dreamlu/deercoder-gin) (单机)  
[micro-go](https://github.com/dreamlu/micro-go) (微服务)  

##### API
- [API 使用](#api-examples)
    - [模型定义](#model)
    - [结构体标记](#struct-gt)
    - [增删改查Crud](#Crud-request)
    - [多表查询](#Crud-More)
    - [批量创建](#createmore)
    - [配置文件模式](#getdevmode)
    - [缓存使用](#cachemanager)
    - [加解密](#aesende)
    - [标准日期](#time)
    - [JSON类型](#jsontype)
    - [字段验证](#validator)  
    - [日志支持](#customlog)
    - [snowflake ID](#snowflakeid)
    - [消息中间件](#msg)
- [扩展 使用](#extend-examples)
    - [crud原生SQL](#crud-selectsql)
    - [更新其他字段](#crud-update)
    - [事务](#transcation)
    - [使用gorm](#use-gorm)
    - [mock假数据](#mock-data)
    - [关于crud的clone](#crud-clone)
    

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
		Inner("order", "user", "order", "service"),
		//Left("order", "service"),
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

#### Struct Gt
> gt:"sub_sql"<忽略该字段解析,<可进行子查询>>  
```go
type Client struct {
	models.AdminCom
	Name    string `gorm:"type:varchar(30)" json:"name" valid:"required,len=2-20"` // 昵称
	Openid  string `json:"openid" gorm:"varchar(30);UNIQUE_INDEX:openid已存在"`       // openID
	HeadImg string `json:"head_img"`                                               // 头像
}

type ClientD struct {
	Client      // 头像
	BuyNum int8 `json:"buy_num" gt:"sub_sql"` // 购买次数
}
// search
func (c *Client) GetBySearch(params cmap.CMap) interface{} {

	buyNumSQL := "(select count(*) from `order` where client_id = `client`.id and status >= 3) as buy_num"
	var datas []*ClientD
	crud.Params(
		gt.Data(&datas),
		gt.SubSQL(buyNumSQL),
		)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}
// output:
print sql: select `id`,`createtime`,`admin_id`,`name`,`openid`,`head_img`,(select count(*) from `order` where client_id = `client`.id and status >= 3) as buy_num from `client` order by id desc limit 0, 2 
```
> gt:"field:fieldName"<特殊情况下,替代json解析>  
```go
    type CVB struct {
		ID          int64      `gorm:"type:bigint(20)" json:"id"`
		ClientVipID int64      `gorm:"type:bigint(20)" json:"client_vip_id"`
		ShopId      int64      `gorm:"type:bigint(20)" json:"shop_id"`
	}

	// 客户行为详情
	type CVBDe struct {
		CVB
		ClientName    string `json:"client_name"`
		VipType       int64  `json:"vip_type" gt:"sub_sql"`
		IsSp          int64  `json:"-" gt:"field:is_sp"`
	}
	gt := &GT{
		Params: &Params{
			InnerTable: []string{"client_vip_behavior", "client_vip", "client_vip", "client"},
			Model:      CVBDe{},
		},
	}
	sqlNt, sql, _, _, _ := GetMoreSearchSQL(gt)
	t.Log(sqlNt)
	t.Log(sql)
// output:
TestGetMoreSearchSQL: db_test.go:267: select count(`client_vip_behavior`.id) as total_num from `client_vip_behavior` inner join `client_vip` on `client_vip_behavior`.`client_vip_id`=`client_vip`.`id`  inner join `client` on `client_vip`.`client_id`=`client`.`id` 
TestGetMoreSearchSQL: db_test.go:268: select `client_vip_behavior`.`id`,`client_vip_behavior`.`client_vip_id`,`client_vip_behavior`.`shop_id`,`client`.`name` as client_name,`client_vip_behavior`.`is_sp` from `client_vip_behavior` inner join `client_vip` on `client_vip_behavior`.`client_vip_id`=`client_vip`.`id`  inner join `client` on `client_vip`.`client_id`=`client`.`id`  order by `client_vip_behavior`.id desc 

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

// id支持批量删除(逗号分割), 如: id = 12,13,14
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

#### Crud More
多表查询支持多表/同一个mysql跨数据库查询/mock假数据等等
```go
    // 多表查询
	// get more search
	var params = make(cmap.CMap)
	//params.Add("user_id", "1")
	//params.Add("key", "梦") // key work
	params.Add("clientPage", "1")
	params.Add("everyPage", "2")
	//params.Add("mock", "1") // mock data
	var or []*OrderD
	crud := NewCrud(
		// 支持同一个mysql多数据库跨库查询
		Inner("gt.order", "user"),
		Left("order", "service"),
		Model(OrderD{}),
		Data(&or),
		//SubWhereSQL("1 = 1", "2 = 2", ""),
	)
	err := crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
	t.Log("\n[User Info]:", or)
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

	err := crud.CreateMore()
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
>  
```go
func TestValidator(t *testing.T) {
    type Test struct {
		ID   int64  `json:"id" gt:"valid:required,min=0,max=5"`
		Name string `json:"name" gt:"valid:required,len=2-5;trans:用户名"`
	}

	// json data
	var test = Test{
		ID:   6,
		Name: "梦",
	}
	t.Log(Valid(test))

	// form data
	var maps = cmap.NewCMap()
	maps["name"] = append(maps["name"], "梦")
	info := ValidForm(maps, Test{})
	//t.Log(info == nil)
	t.Log(info)
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

### Msg  
消息中间件nsq  
> 前提: 运行nsq, 参考[docker-compose.yaml](./test/docker/nsq/docker-compose.yaml)  

参考:[nsq_test.go](./msg/nsq_test.go)  

### Extend Examples  

#### Crud Selectsql
```go
    sql := "update `user` set name=? where id=?"
    t.Log("[Info]:", crud.Select(sql, "梦sql", 1).Exec())
    var user []*User
	cd := crud.Params(
		Data(&user),
		//ClientPage(1), // serch()分页需要
		//EveryPage(2),
	).
		Select("select *from user").
		Select("where id > 0")
	if true {
		cd.Select("and 1=1")
	}
	cd.Search() // 查询 + 分页
    t.Log(crud.Pager()) // Search()分页数据
	//cd.Single() // 注释Search使用Single()直接查询
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

#### Use Gorm
> 模型定义需遵循:  
> gt v1.20以前参考[模型定义v1](https://v1.gorm.io/zh_CN/docs/models.html)  
> gt v1.20+参考[模型定义v2](https://gorm.io/zh_CN/docs/models.html)
```go
// example 1:
// 根据模型定义自动生成表
gt.NewDBTool().AutoMigrate(&User{},&Order{})
// 直接使用gorm:
db := gt.NewCrud().DB()
```

#### Mock Data
> 使用mock参数, 生成随机数据, 将不会进行数据库查询  
```go
    GetBySearch(params cmap.CMap) Crud     // search
    GetByData(params cmap.CMap) Crud       // get data no search
    GetMoreBySearch(params cmap.CMap) Crud // more search    
    // 以上三种支持mock参数,传递的参数mock=1即可
```  
ps:  
1.不支持CJSON类型, 请使用tag: `faker:"-"`进行过滤  
2.不支持图片等实体文件数据  
3.默认随机生成,如有长度等其他要求,请参考:[faker_test](https://github.com/bxcodec/faker/blob/master/faker_test.go)  

#### Crud clone
```go
// 关于crud的
一个crud = gt.NewCrud()对象中的参数是共享的,通用的增删改查针对同一张表可复用  
如果进行了表关联或改变了模型, 需要重新cd = gt.NewCrud(),否则继续使用crud容易影响到其他使用这个变量的地方
```

- 约定  
1.模型结构体json 内容与表字段保持一致  
2.返回格式参考[result](tool/result/result.go)    
3.多表关联命名, 模型中其他表字段命名: `他表名 + "_" + 他表字段名`  
n....  