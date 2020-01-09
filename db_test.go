// package gt

package gt

import (
	"fmt"
	"github.com/dreamlu/go-tool/tool/type/json"
	"github.com/dreamlu/go-tool/tool/type/time"
	"log"
	"testing"
	time2 "time"
)

// params map[string][]string is web request GET params
// in golang, it was url.values

type User struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Createtime time.CTime `json:"createtime"`
}

type UserInfo struct {
	ID       int64      `json:"id"`
	UserID   int64      `json:"user_id"`   //用户id
	UserName string     `json:"user_name"` //用户名
	Userinfo json.CJSON `json:"userinfo"`
}

// order
type Order struct {
	ID         int64 `json:"id"`
	UserID     int64 `json:"user_id"`     // user id
	ServiceID  int64 `json:"service_id"`  // service table id
	CreateTime int64 `json:"create_time"` // createtime
}

// order detail
type OrderD struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`      // user id
	UserName    string     `json:"user_name"`    // user table column name
	ServiceID   int64      `json:"service_id"`   // service table id
	ServiceName string     `json:"service_name"` // service table column `name`
	Createtime  time.CTime `json:"createtime"`   // createtime
}

// 局部
var crud = NewCrud()

func TestDB(t *testing.T) {

	var user = User{
		Name: "测试xx",
		//Createtime:JsonDate(time.Now()),
	}

	// return create id
	_ = crud.DB().CreateData(&user)
	t.Log("user: ", user)
	t.Log(crud.DB().RowsAffected)
	user.Name = "haha"
	_ = crud.DB().CreateData(&user)
	t.Log("user: ", user)
	t.Log(crud.DB().RowsAffected)
	//user.ID = 8 //0
	//ss = UpdateStructData(&user)
	//log.Println(ss)
}

// 通用分页测试
// 如：
func TestSqlSearch(t *testing.T) {
	sql := fmt.Sprintf(`select a.id,a.user_id,a.userinfo,b.name as user_name from userinfo a inner join user b on a.user_id=b.id where 1=1 and `)
	sqlNt := `select
		count(distinct a.id) as total_num
		from userinfo a inner join user b on a.user_id=b.id
		where 1=1 and `
	var ui []UserInfo

	//页码,每页数量
	clientPage := int64(1) //默认第1页
	everyPage := int64(10) //默认10页

	//可定制
	//args map[string][]string
	//look go-tool/demo
	//args is url.values
	//for k, v := range args {
	//	if k == "clientPage" {
	//		clientPageStr = v[0]
	//		continue
	//	}
	//	if k == "everyPage" {
	//		everyPageStr = v[0]
	//		continue
	//	}
	//	if v[0] == "" { //条件值为空,舍弃
	//		continue
	//	}
	//	v[0] = strings.Replace(v[0], "'", "\\'", -1) //转义
	//	sql += "a." + k + " = '" + v[0] + "' and "
	//	sqlNt += "a." + k + " = '" + v[0] + "' and "
	//}

	sql = string([]byte(sql)[:len(sql)-4]) //去and
	sqlNt = string([]byte(sqlNt)[:len(sqlNt)-4])
	sql += "order by a.id "
	log.Println(crud.DB().GetDataBySQLSearch(&ui, sql, sqlNt, clientPage, everyPage, nil, nil))
	log.Println(ui[0].Userinfo.String())
}

// 常用分页测试(两张表)
// 如:
func TestSqlSearchV2(t *testing.T) {
	//var ui []UserInfo
	//
	////args map[string][]string
	////look github.com/dreamlu/deercoder-gin
	////args is url.values
	//log.Println(GetDoubleTableDataBySearch(UserInfo{},&ui, "userinfo", "user", args))
	//log.Println(ui)
}

// select 数据存在验证
func TestValidateData(t *testing.T) {
	sql := "select *from `user` where id=2"
	ss := crud.DB().ValidateSQL(sql)
	log.Println(ss)
}

// 分页搜索中key测试
func TestGetSearchSql(t *testing.T) {

	type UserDe struct {
		User
		Num int64 `json:"num" gt:"sub_sql"`
	}

	var args = make(map[string][]string)
	args["clientPage"] = append(args["clientPage"], "1")
	args["everyPage"] = append(args["everyPage"], "2")
	//args["key"] = append(args["key"], "梦 嘿,伙计")
	sub_sql := ",(select aa.name from shop aa where aa.user_id = a.id) as shop_name"
	sqlNt, sql, _, _, _ := GetSearchSQL(&GT{
		Params: args,
		Table:  "user",
		Model:  UserDe{},
		SubSQL: sub_sql,
	})
	log.Println("SQLNOLIMIT:", sqlNt, "\nSQL:", sql)

	// 两张表，待重新测试
	// sqlNt, sql, _, _ = GetDoubleSearchSQL(UserInfo{}, "userinfo", "user", args)
	// log.Println("SQLNOLIMIT==>2:", sqlNt, "\nSQL==>2:", sql)

}

// 通用sql以及参数
func TestGetDataBySql(t *testing.T) {
	var sql = "select id,name,createtime from `user` where id = ?"

	var user User
	_ = crud.DB().GetDataBySQL(&user, sql, "1")
	log.Println(user)

	//GOTool.Raw(sql, []interface{}{1, "梦"}[:]...).Scan(&user)
	//log.Println(user)
}

func TestGetDataBySearch(t *testing.T) {
	var args = make(map[string][]string)
	args["name"] = append(args["name"], "梦")
	args["key"] = append(args["key"], "梦")
	args["clientPage"] = append(args["clientPage"], "1")
	args["everyPage"] = append(args["everyPage"], "2")
	var user []*User
	_, _ = crud.DB().GetDataBySearch(&GT{
		Params: args,
		Table:  "user",
		Model:  User{},
		Data:   &user,
	})
	t.Log(user[0])
}

// 通用增删该查测试
func TestCrud(t *testing.T) {
	var args = make(map[string][]string)
	args["name"] = append(args["name"], "梦")

	// var crud DbCrud
	// must use AutoMigrate
	// get by id
	//GOTool.DB.AutoMigrate(&User{})
	//crud.DB().DB.AutoMigrate(&User{})
	var user User
	//param := &Params{
	//	Table:     "user",
	//	Data: &user,
	//}
	crud = NewCrud(
		Table("user"),
		Data(&user),
	)
	info := crud.GetByID("1")
	log.Println(info, "\n[User Info]:", user)

	// get by search
	var users []*User
	//param = &Params{
	//	Table:     "user",
	//	Model:     User{},
	//	Data: &users,
	//}
	//crud = NewCrud(param)
	crud.Params(
		Model(User{}),
		Data(&user),
	)
	args["name"][0] = "梦"
	crud.GetBySearch(args)
	log.Println("\n[User Info]:", users)

	// delete
	info2 := crud.Delete("12")
	log.Println(info2)

	// update
	args["id"] = append(args["id"], "4")
	args["name"][0] = "梦4"
	info2 = crud.UpdateForm(args)
	log.Println(info2)

	// create
	//var args2 = make(map[string][]string)
	//args2["name"] = append(args2["name"],"梦c")
	////db  = DbCrud{"user", nil,&user}
	//info = db.Create(args2)
	//log.Println(info)

}

// 通用增删改查sql测试
func TestCrudSQL(t *testing.T) {
	//var db = DbCrud{}
	sql := "update `user` set name=? where id=?"
	log.Println("[Info]:", crud.Select(sql, "梦sql", 1).Exec())
	log.Println("[Info]:", crud.DB().RowsAffected)
}

// 测试多表连接
func TestGetMoreDataBySearch(t *testing.T) {
	// 多表查询
	// get more search
	var params = make(map[string][]string)
	params["user_id"] = append(params["user_id"], "1")
	params["key"] = append(params["key"], "梦")
	params["clientPage"] = append(params["clientPage"], "1")
	params["everyPage"] = append(params["everyPage"], "2")
	var or []*OrderD
	//param := &Params{
	//	InnerTable: []string{"order", "user"},
	//	LeftTable:  []string{"service"},
	//	Model:      OrderD{},
	//	Data:  &or,
	//}
	crud := NewCrud(
		InnerTable([]string{"order", "user", "order", "service"}),
		//LeftTable([]string{"order", "service"}),
		Model(OrderD{}),
		Data(&or),
	)
	_, err := crud.GetMoreBySearch(params)
	if err != nil {
		log.Println(err)
	}
	t.Log("\n[User Info]:", or[0])
}

func TestGetMoreSearchSQL(t *testing.T) {
	type ClientVipBehavior struct {
		ID          int64      `gorm:"type:bigint(20)" json:"id"`
		ClientVipID int64      `gorm:"type:bigint(20)" json:"client_vip_id"`
		ShopId      int64      `gorm:"type:bigint(20)" json:"shop_id"`
		StaffId     int64      `gorm:"type:bigint(20)" json:"staff_id"`
		Status      int64      `gorm:"type:tinyint(2);DEFAULT:0" json:"status"`
		Num         int64      `json:"num" gorm:"type:int(11)"` // 第几次参加
		Createtime  time.CTime `gorm:"type:datetime" json:"createtime"`
	}

	// 客户行为详情
	type ClientVipBehaviorDe struct {
		ClientVipBehavior
		ClientName    string `json:"client_name"`
		ClientHeadimg string `json:"client_headimg"`
		VipType       int64  `json:"vip_type" gt:"sub_sql"` // 0意向会员, 1会员
		IsSp          int64  `json:"is_sp" gt:"sub_sql"`    // 是否代言人, 0不是, 1是
	}
	gt := &GT{
		InnerTable: []string{"client_vip_behavior", "client_vip", "client_vip", "client"},
		Model:      ClientVipBehaviorDe{},
	}
	_, sql, _, _, _ := GetMoreSearchSQL(gt)
	log.Println(sql)
}

// 批量创建
func TestCreateMoreData(t *testing.T) {

	var user = []User{
		{Name: "测试1", Createtime: time.CTime(time2.Now())},
		{Name: "测试2"},
	}

	//param := &Params{
	//	Table: "user",
	//	Model: User{},
	//}
	crud := NewCrud(
		Table("user"),
		Model(User{}),
		Data(user),
		//SubSQL("(asdf) as a","(asdfa) as b"),
	)

	err := crud.CreateMoreData()
	log.Println(err)
}

// 继承tag解析测试
func TestExtends(t *testing.T) {
	type UserDe struct {
		User
		Other string `json:"other"`
	}

	type UserDeX struct {
		a []string
		UserDe
		OtherX string `json:"other_x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name"`
		UserDeX
	}
	t.Log(GetColSQL(UserDeX{}))
	t.Log(GetMoreTableColumnSQL(UserMore{}, []string{"user", "shop"}[:]...))
}

// select test
func TestDBCrud_Select(t *testing.T) {
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
	crud.Search()
	crud.Single()
}

// test update/delete
func TestDBCrud_Update(t *testing.T) {

	crud := crud.Params(
		Table("user"),
		Data(&User{
			ID:   1,
			Name: "梦S",
		}),
	)
	crud.Update()
	t.Log(crud.DB().RowsAffected)
	crud.Params(Data(&User{
		ID:   1,
		Name: "梦SSS",
	}))
	crud.Create()
	t.Log(crud.DB().RowsAffected)
}

// test update/delete
func TestDBCrud_Create(t *testing.T) {

	crud.Params(
		Table("user"),
		Data(&User{
			ID:1,
			Name: "梦S",
		}),
	).Create()
	t.Log(crud.DB().Error)
	crud.Params(Data(&User{
		Name: "梦SSS",
	})).Create()
	t.Log(crud.DB().Error)
}
