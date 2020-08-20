package gt

import (
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type Client struct {
	ID         primitive.ObjectID `json:"id"`
	Name       string             `json:"name" gt:"valid:len=3-5;trans:名称"`
	BirthDate  time.CDate         `json:"birth_date" gorm:"type:date"` // data
	CreateTime time.CTime         `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"create_time"`
	Account    float64            `json:"account" gorm:"type:decimal(10,2)"`
}

func TestMongo_Create(t *testing.T) {

	var user = Client{
		Name: "test",
	}
	cd := NewCrud(
		D("mongo"),
		Model(Client{}),
		Data(user),
	)
	cd.Create()
	cd.Params(Data([]Client{user, user}))
	cd.CreateMore()
}

func TestMongo_Update(t *testing.T) {

	var user = Client{
		//ID:         "5f3cd15cf2f80b74c05f5033",
		Name:       "",
		BirthDate:  time.CDate{},
		CreateTime: time.CTime{},
		Account:    0,
	}
	cd := NewCrud(
		D("mongo"),
		Table("client"),
		Data(user),
	)
	//filter := bson.D{{"_id", "Ash"}}
	//
	//update := bson.D{
	//	{"$inc", bson.D{
	//		{"age", 1},
	//	}},
	//}
	cd.Update()
}

func TestMongo_GetByData(t *testing.T) {

	var client []*Client
	cd := NewCrud(
		D("mongo"),
		Table("client"),
		Data(&client),
	)
	cd.GetByData(cmap.NewCMap().Set("name", "update"))
	t.Log(client)

	var client2 Client
	cd = NewCrud(
		D("mongo"),
		Table("client"),
		Data(&client2),
	)
	cd.GetByData(cmap.NewCMap().Set("name", "update"))
	t.Log(client2)
}

func TestMongo_GetByID(t *testing.T) {

	var client Client
	cd := NewCrud(
		D("mongo"),
		Table("client"),
		Data(&client),
	)
	objID, _ := primitive.ObjectIDFromHex("5f3cd15cf2f80b74c05f5033")
	t.Log(objID.String())
	cd.GetByID(objID)
	t.Log(client)
}

func TestMongo_GetBySearch(t *testing.T) {

	var client []*Client
	cd := NewCrud(
		D("mongo"),
		Table("client"),
		Data(&client),
	)
	cd.GetBySearch(cmap.NewCMap().Set("clientPage", "1").Set("everyPage", "3"))
	t.Log(cd.Error())
	t.Log(client)
}

func TestMongo_Delete(t *testing.T) {

	var client []*Client
	cd := NewCrud(
		D("mongo"),
		Model(Client{}),
		Data(&client),
	)
	objID, _ := primitive.ObjectIDFromHex("5f3e64d46a099e10d7879f64")
	cd.Delete(objID)
	t.Log(cd.Error())
	t.Log(cd.RowsAffected())
}
