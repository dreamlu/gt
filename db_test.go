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

// 全局
//var GOTool = NewDBTool()
// 局部
var crud = NewCrud()

func TestDB(t *testing.T) {

	var user = User{
		Name: "测试xx",
		//Createtime:JsonDate(time.Now()),
	}

	//err := CreateDataJ(&user)
	//log.Println(err)

	// return create id
	_ = crud.DB().CreateDataJ(&user)
	log.Println("user: ", user)

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
	log.Println(crud.DB().GetDataBySQLSearch(&ui, sql, sqlNt, clientPage, everyPage))
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
		Params:    args,
		Table:     "user",
		Model:     User{},
		ModelData: &user,
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
	//	ModelData: &user,
	//}
	crud = NewCrud(
		Table("user"),
		ModelData(&user),
	)
	info := crud.GetByID("1")
	log.Println(info, "\n[User Info]:", user)

	// get by search
	var users []*User
	//param = &Params{
	//	Table:     "user",
	//	Model:     User{},
	//	ModelData: &users,
	//}
	//crud = NewCrud(param)
	crud.Params().Model = User{}
	crud.Params().ModelData = &users
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
	log.Println("[Info]:", crud.UpdateBySQL(sql, "梦sql", 1))
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
	//	ModelData:  &or,
	//}
	crud := NewCrud(
		InnerTable([]string{"order", "user"}),
		LeftTable([]string{"service"}),
		Model(OrderD{}),
		ModelData(&or),
	)
	_, err := crud.GetMoreBySearch(params)
	if err != nil {
		log.Println(err)
	}
	t.Log("\n[User Info]:", or[0])
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
		//SubSQL("(asdf) as a","(asdfa) as b"),
	)

	err := crud.CreateMoreData(user)
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
