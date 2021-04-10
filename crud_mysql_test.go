// package gt

package gt

import (
	sql2 "database/sql"
	json2 "encoding/json"
	"fmt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"github.com/dreamlu/gt/tool/type/time"
	"log"
	"net/http"
	"runtime"
	"testing"
	time2 "time"
)

// params param.CMap is web request GET params
// in golang, it was url.values

// user model
type User struct {
	ID         uint64     `json:"id"`
	Name       string     `json:"name" gt:"valid:len=3-5;trans:名称" gorm:"<-:update"`
	BirthDate  time.CDate `gorm:"type:date"` // data
	CreateTime time.CTime `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"create_time"`
	Account    float64    `json:"-" gorm:"type:decimal(10,2)"`
}

type UserInfo struct {
	ID   uint64 `json:"id"`
	Some string `json:"some"`
}

func (u User) String() string {
	b, _ := json2.Marshal(u)
	return string(b)
}

// service model
type Service struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// order model
type Order struct {
	ID         uint64      `json:"id"`
	UserID     int64       `json:"user_id"` // user id
	UserInfoID uint64      `json:"user_info_id"`
	ServiceID  int64       `json:"service_id"` // service table id
	CreateTime time.CTime  `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"create_time"`
	StartTime  time.CSTime `json:"start_time"`
	EndTime    time.CSTime `json:"end_time"`
}

// order detail
type OrderD struct {
	Order
	UserName    string     `json:"user_name" gt:"field:user.name"`  // user table column name
	ServiceName string     `json:"service_name"`                    // service table column `name`
	Info        json.CJSON `json:"info" gt:"sub_sql" faker:"cjson"` // json
	BirthDate   time.CDate `gorm:"type:date"`                       // data
}

// 局部
var crud Crud

func init() {
	err := DB().AutoMigrate(User{}, Order{}, Service{}, UserInfo{})
	fmt.Println(err)
	crud = NewCrud()
}

func TestDB(t *testing.T) {

	var user = User{
		Name:      "测试xx",
		BirthDate: time.CDate(time2.Now()),
		//Createtime:JsonDate(time.Now()),
	}

	// return create id
	DB().Create("", &user)
	t.Log("user: ", user)
	user.Name = "haha"
	DB().Create("", &user)
	t.Log("user: ", user)
	var user2 User
	crud.Params(Model(User{}), Data(&user2)).GetByID(1)
	t.Log(user2)
}

// 通用增删该查测试
// 传参可使用url.Values替代param.CMap操作方便
func TestCrud(t *testing.T) {

	// add
	var user = User{
		Name: "new name",
	}
	crud = NewCrud()
	cd := crud.Params(
		Model(User{}),
		Table(""),
		Data(&user),
	).Create()

	t.Log(cd.Error())
	// update
	user.ID = 2
	crud.Update()

	// get by id
	var user2 User
	crud.Params(Data(&user2)).GetByID(2)
	t.Log(user2, "\n[GetByID]:", crud.Error())

	// get by search
	//var args = url.Values{}
	//args.Add("name", "梦")
	// get by search
	var users []*User
	crud.Params(
		Model(User{}),
		Data(&users),
		WhereSQL("1=1"),
	)
	crud = crud.Params(Table("gt.user")).GetBySearch(cmap.NewCMap().Set("id", "1000"))
	t.Log("\n[User Info]:", users)
	t.Log(crud.Error())

	// delete
	info2 := crud.Delete(12)
	t.Log(info2.Error())
	info2 = crud.Delete("12,13,14")
	t.Log(info2.Error())

	// update by form request
	//args.Add("id", "4")
	//args.Set("name", "梦4")
	//err := crud.UpdateForm(cmap.CMap(args))
	//log.Println(err)
}

// select sql
func TestCrudSQL(t *testing.T) {
	crud := NewCrud()
	sql := "update `user` set name=? where id = ?"
	t.Log("[Info]:", crud.Select(sql, "梦sql", 1).Select("and 1=1 and").
		Select(&User{
			Name: "梦S",
		}).Exec())
	t.Log("[Info]:", crud.Select(sql, "梦sql", 1).Exec())
	t.Log("[Info]:", crud.RowsAffected())
	var user []User
	sql = "select * from `user` where name=? and id=?"
	cd := NewCrud()
	t.Log("[Info]:", cd.Params(Data(&user)).Select(sql, "梦sql", 1).Select(" and 1=1").Exec())
	t.Log("[Info]:", cd.Params(Data(&user)).Select(sql, "梦sql", 1).Exec())
}

func TestGetSearchSql(t *testing.T) {

	type UserDe struct {
		User
		Num int64 `json:"num" gt:"sub_sql"`
	}

	var args = make(cmap.CMap)
	args.Add("clientPage", "1")
	args.Add("everyPage", "2")
	//args["key"] = append(args["key"], "梦 嘿,伙计")
	//sub_sql := ",(select aa.name from shop aa where aa.user_id = a.id) as shop_name"
	gt := &GT{
		Params: &Params{
			Table: "user",
			Model: User{},
		},
		CMaps:  args,
		Select: "",
		From:   "",
		Group:  "",
		Args:   nil,
	}
	gt.GetSearchSQL()
	t.Log("SQLNOLIMIT:", gt.sqlNt, "\nSQL:", gt.sql)
}

func TestGetDataBySql(t *testing.T) {
	var sql = "select id,name,create_time from `user` where id = ?"

	var user User
	err := crud.Params(Data(&user)).Select(sql, "1000").Single().Error()
	t.Log(err)
	t.Log(user)
}

func TestGetDataBySearch(t *testing.T) {
	var args = make(cmap.CMap)
	args["key"] = append(args["key"], "梦")
	args["clientPage"] = append(args["clientPage"], "1")
	args["everyPage"] = append(args["everyPage"], "2")
	var user []*User
	DB().GetBySearch(&GT{
		CMaps: args,
		Params: &Params{
			Table: "user",
			Model: User{},
			Data:  &user,
		},
	})
	t.Log(DB().res.Error)
	if len(user) > 0 {
		t.Log(user[0])
	}
}

// 测试多表连接
func TestGetMoreDataBySearch(t *testing.T) {

	type Key struct {
		UserName    string `json:"user_name"`
		UserAccount string `json:"user_account"`
	}
	// 多表查询
	// get more search
	var params = cmap.NewCMap().
		Set("key", "test 1"). // key work
		Set("clientPage", "1").
		Set("everyPage", "2")
	//params.Add("mock", "1") // mock data
	var or []*OrderD
	crud := NewCrud(
		// 支持同一个mysql多数据库跨库查询
		Inner("order", "gt.user"),
		Left("order", "service"),
		Model(OrderD{}),
		Data(&or),
		//KeyModel(Key{}),
		WhereSQL("1 = ?", 1).WhereSQL("2 = ?", 2),
		Distinct("order.id"),
	)
	err := crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
	err = crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
	for _, v := range or {
		t.Log(v)
	}
}

func TestGetMoreSearchSQL(t *testing.T) {
	type CVB struct {
		ID          int64 `gorm:"type:bigint(20)" json:"id"`
		ClientVipID int64 `gorm:"type:bigint(20)" json:"client_vip_id"`
		ShopId      int64 `gorm:"type:bigint(20)" json:"shop_id"`
	}

	// 客户行为详情
	type CVBDe struct {
		CVB
		ClientName string `json:"client_name"`
		VipType    int64  `json:"vip_type" gt:"sub_sql"`
		IsSp       int64  `json:"-" gt:"field:is_sp"`
	}
	gt := &GT{
		Params: &Params{
			InnerTable: []string{"client_vip_behavior", "client_vip", "client_vip", "client"},
			Model:      CVBDe{},
		},
	}
	gt.GetMoreSQL()
	t.Log(gt.sqlNt)
	t.Log(gt.sql)
}

// 批量创建
func TestCreateMoreData(t *testing.T) {

	type Par struct {
		Account float64 `json:"account"`
	}

	type UserPar struct {
		Name       string     `json:"name" gt:"valid:len=2-10"`
		CreateTime time.CTime `json:"create_time"`
		Par
	}
	type User struct {
		ID uint64 `json:"id"`
		UserPar
	}
	type UserParP struct {
		UserPar
	}

	var up = []UserParP{
		{
			UserPar{Name: "测试1", CreateTime: time.CTime(time2.Now())},
		},
		{
			UserPar{Name: "测试2"},
		},
	}
	crud := NewCrud(
		Table("user"),
		Model(UserParP{}),
		Data(up),
		//SubSQL("(asdf) as a","(asdfa) as b"),
	)

	err := crud.CreateMore().Error()
	t.Log(err)
}

// 继承tag解析测试
func TestExtends(t *testing.T) {
	type UserDe struct {
		User
		Other string `json:"other"`
	}

	type UserDeX struct {
		a []string `gt:"ignore"`
		UserDe
		OtherX string `json:"other_x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name"`
		UserDeX
	}
	for i := 0; i < 3; i++ {
		t.Log(GetColSQL(UserDeX{}))
		t.Log(GetMoreTableColumnSQL(UserMore{}, []string{"user", "shop"}[:]...))
	}
}

// select test
func TestDBCrud_Select(t *testing.T) {
	var (
		params = cmap.NewCMap()
		user   []*User
	)
	params.Add("clientPage", "1")
	params.Add("everyPage", "2")
	//params.Add("mock", "1")
	cd := crud.Params(
		Data(&user),
	).
		Select("select *from user").
		Select("where id > 0")
	if true {
		cd.Select("and 1=1")
	}
	// search
	cd = cd.Search(params)
	t.Log(cd.Error())
	t.Log(user)
	// single
	cd2 := crud.Params(
		Data(&user),
	).
		Select("select *from user limit 2")
	t.Log(cd2.Single().Error())
	_, file, line, ok := runtime.Caller(1)
	if ok {
		t.Log(file, "[]", line)
	}

	// use gorm 2.0 support basic type scan replace
	var name sql2.NullString
	cd3 := NewCrud(Data(&name)).Select("select ifnull(name,'') from user where id = 10").Single()
	t.Log(cd3.Error())
	t.Log(name)

	var names []string
	cd4 := NewCrud(Data(&names)).Select("select ifnull(name,'') from user where id > 0 limit 2").Single()
	t.Log(cd4.Error())
	t.Log(names)
}

// test update/delete
func TestDBCrud_Update(t *testing.T) {

	crud := crud.Params(
		//Table("user"),
		Model(User{}),
		Data(&User{
			ID:   1,
			Name: "梦sql",
		}),
	)
	cd := crud.Update()
	t.Log(cd.Error(), cd.RowsAffected())
	t.Log(crud.Params(Data(User{
		ID:   1,
		Name: "梦sql",
	})).Select("name = ?", "梦sql").Update().RowsAffected())
	t.Log(crud.Error())
}

// test update/delete
func TestDBCrud_Create(t *testing.T) {

	crud.Params(
		Table("user"),
		Data(&User{
			ID:   11234,
			Name: "梦S",
		}),
	)
	t.Log(crud.Error())
	t.Log(crud.Create().Error())
	crud.Params(
		Data(&User{
			Name: "梦SSS2",
		})).Create()
	t.Log(crud.Error())
}

func TestGetColSQLAlias(t *testing.T) {
	sql := GetColSQLAlias(User{}, "a")
	t.Log(sql)

	sql = GetColSQLAlias(OrderD{}, "a")
	t.Log(sql)
}

func httpServerDemo(w http.ResponseWriter, r *http.Request) {
	// get more search
	var params = make(cmap.CMap)
	params.Add("mock", "1") // mock data
	var or []*OrderD
	crud := NewCrud(
		Inner("order", "user"),
		//Left("order", "service"),
		Model(OrderD{}),
		Data(&or),
		SubWhereSQL("1 = 1", "2 = 2", ""),
	)
	err := crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
	_, _ = fmt.Fprintf(w, result.GetSuccess(or).String())
}

// test mock data
func TestMock(t *testing.T) {
	//http.HandleFunc("/", httpServerDemo)
	//log.Println("http://127.0.0.1:9090")
	//err := http.ListenAndServe(":9090", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
}

func TestField(t *testing.T) {
	type Sv struct {
		ID   uint64 `json:"id"`
		Name string `json:"name" gorm:"type:varchar(50)"` // 名称
	}

	// 详情
	type SvD struct {
		Sv
		GoodsID   uint64 `json:"goods_id" gt:"field:goods.id"` // 商品id
		GoodsName string `json:"goods_name"`                   // 商品名
	}

	var datas []*SvD
	crud.Params(
		Model(SvD{}),
		Data(&datas),
		Inner("sv.sv:id", "sv.svg:sv_id"),
		Left("sv.svg", "shop.goods"),
	)
	var params = cmap.NewCMap()
	_ = crud.GetMoreBySearch(params)
}

// test transaction
func TestTransaction(t *testing.T) {
	cd := NewCrud().Begin()
	type UserP struct {
		Name string
	}
	var (
		params = cmap.NewCMap()
		userP  = UserP{
			Name: "test1",
		}
		user = User{
			ID:   1,
			Name: "test2",
		}
		users []*User
	)
	_ = cd.Params(Data(&user)).Select("select *from user where id = 1").Single()
	t.Log("step1: ", user)

	user.Name = "testUpdate"
	cd.Params(Data(user)).Update()

	user.Name = "testUpdate"
	user.ID = 0
	cd.Params(Table("user"), Data(userP)).Select("name = ?", "梦").Update()

	cd.SavePoint("point1")

	params.Set("id", "1").Set("name", "sql")
	cd.Params(Data(&user)).Get(params)
	t.Log("step2: ", user)

	cd.Params(Data(&users)).GetBySearch(params)
	for _, v := range users {
		t.Log("step3: ", v)
		break
	}

	cd.Select("update user set name = 'testExec' where id = 1").Exec()
	cd.Params(Data(&user)).GetByID(1)
	t.Log("step4: ", user)

	cd.RollbackTo("point1")
	params.Set("id", "1")
	cd.Params(Data(&user)).Get(params)
	t.Log("point1: ", user)

	err := cd.Params(Data(user)).Create().Error()
	t.Log("error", err)

	cd.Params(Data(&user)).GetByID(2)
	t.Log("step5: ", user)

	cd.Params(Data(user)).Update()

	cd.Update()

	cd.Create()

	user.ID = 1
	user.Name = "test3"
	cd.Params(Data(user)).Update()

	cd.Commit()
}

func TestDBDouble10to2(t *testing.T) {
	var user User
	cd := NewCrud(
		Data(&user),
		Model(User{}),
	).
		Get(cmap.NewCMap().Set("order", "id desc"))
	t.Log(user)
	t.Log(cd.Error())

	var user2 User
	NewCrud(
		Data(&user2),
		Model(User{}),
	).Select("select *from `user` where id = ?", 1).Single()
	t.Log(user2)
}

func TestMysql_GetMoreByData(t *testing.T) {
	var or []*OrderD
	cd := NewCrud(
		// 支持同一个mysql多数据库跨库查询
		Inner("order", "gt.user"),
		Left("order", "service"),
		Model(OrderD{}),
		Data(&or),
		KeyModel(OrderD{}),
		//SubWhereSQL("1 = 1", "2 = 2", ""),
	)
	err := cd.GetMore(cmap.NewCMap()).Error()
	if err != nil {
		t.Error(err)
	}
	for _, v := range or {
		t.Log(v)
	}
}

// 测试mysql仅时分秒格式
func TestMysql_Time(t *testing.T) {
	var or []*Order
	cd := NewCrud(
		Model(Order{}),
		Data(&or),
	)
	err := cd.Get(cmap.NewCMap()).Error()
	if err != nil {
		t.Error(err)
	}
	for _, v := range or {
		t.Log(v.EndTime.String())
		s := time.CTime(v.StartTime)
		t.Log(s.String())
		t.Log(v)
	}

	cd.Params(Data(&Order{
		StartTime: time.ParseCSTime("12:00:00"),
		EndTime:   time.ParseCSTime("12:00:00"),
	})).Create()
}

// select * unique table column test
func TestGetMoreSearchResolve(t *testing.T) {

	// order detail
	type OrderD struct {
		Order
		UserBirthDate string `json:"user_birth_date"` // user table column name
		Name          string `json:"name"`            // user table column name
	}
	var params = cmap.NewCMap().
		Set("key", "test 1"). // key work
		Set("clientPage", "1").
		Set("everyPage", "2")
	var or []*OrderD
	crud := NewCrud(
		Inner("order:user_id", "user:id"),
		Model(OrderD{}),
		Data(&or),
		//KeyModel(Key{}),
	)
	err := crud.GetMoreBySearch(params).Error()
	if err != nil {
		log.Println(err)
	}
}

func TestMysql_GetMoreBySearchInnerLeftCondition(t *testing.T) {
	type UserAndUserInfo struct {
		UserName     string // will user.name
		UserInfoSome string `gt:"field:user_info.some"`
	}
	var data []*UserAndUserInfo
	crud.Params(
		Data(&data),
		Model(UserAndUserInfo{}),
		Inner("order:user_id,user_id,id=1", "user:id,id,id=1,id=2"), // "order", "user_info"), // inner/left join on a.column = b.column and ...
		Left("order:user_id", "user_info:id,id=1"),
	)
	var params = make(cmap.CMap)
	crud.GetMoreBySearch(params)
}
