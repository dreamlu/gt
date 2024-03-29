package mongo

import (
	"github.com/dreamlu/gt/src/type/cmap"
	"github.com/dreamlu/gt/src/type/time"
	"github.com/dreamlu/gt/src/type/tmap"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"testing"
	time2 "time"
)

type Client struct {
	ID         string  `json:"id" bson:"_id,omitempty"`
	Name       string  `json:"name" gt:"valid:len=3-5;trans:名称"`
	BirthDate  string  `json:"birth_date" gorm:"type:date"` // data
	CreateTime string  `json:"create_time"`
	Account    float64 `json:"account" gorm:"type:decimal(10,2)"`
}

func TestMongo_Create(t *testing.T) {

	var user = Client{
		Name:       "test",
		CreateTime: time.CTime(time2.Now()).String(),
	}
	cd := NewCrud(
		Model(Client{}),
		Data(user),
	)
	cd.Create()
	cd.Params(Data([]Client{user, user}))
	cd.CreateMore()
}

func TestMongo_Update(t *testing.T) {

	var user = Client{
		ID:         "5f4372d5b7f7ce9d6e6ba479",
		Name:       "new_update",
		CreateTime: time.CTime(time2.Now()).String(),
		Account:    0,
	}
	cd := NewCrud(
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
	t.Log(cd.Error())
}

func TestMongo_Get(t *testing.T) {

	var client []*Client
	cd := NewCrud(

		Table("client"),
		Data(&client),
	)
	cd.Find(cmap.NewCMap().Set("name", "update"))
	t.Log(client)

	var client2 Client
	cd = NewCrud(

		Table("client"),
		Data(&client2),
	)
	cd.Find(cmap.NewCMap().Set("name", "update"))
	t.Log(client2)
}

func TestMongo_GetByID(t *testing.T) {

	var client Client
	cd := NewCrud(
		Table("client"),
		Data(&client),
	)
	objID, _ := primitive.ObjectIDFromHex("609b451acb2ae879ea3fe8e9")
	t.Log(objID.String())
	cd.FindID(objID)
	t.Log(client)
}

func TestMongo_GetBySearch(t *testing.T) {

	var client []*Client
	cd := NewCrud(
		Model(Client{}),
		//Table("client"),
		Data(&client),
	)
	cd.FindSearch(tmap.NewTMap[string, any]().
		//Set("clientPage", "1").
		//Set("everyPage", "3").
		Set("order", "id desc").
		//Set("name", "test"),
		Set("account", 0),
	//Set("create_time", "2020-08-24 16:03:55"),
	)
	t.Log(cd.Error())
	for _, v := range client {
		t.Log(v)
	}
	t.Log(cd.Pager())
}

func TestMongo_Delete(t *testing.T) {

	var client []*Client
	cd := NewCrud(

		Model(Client{}),
		Data(&client),
	)
	objID, _ := primitive.ObjectIDFromHex("5f3e64d46a099e10d7879f64")
	cd.Delete(objID)
	t.Log(cd.Error())
	t.Log(cd.RowsAffected())
}

func TestMongo_GoRoutine(t *testing.T) {

	cd := NewCrud(
		Model(Client{}),
	)

	var g sync.WaitGroup
	for i := 0; i < 100; i++ {
		g.Add(1)
		go func() {
			defer g.Done()
			var client []*Client
			cd.Params(Data(&client))
			cd.FindSearch(tmap.NewTMap[string, any]().
				Set("order", "id desc"),
			)
			t.Log(len(client))
		}()
	}
	g.Wait()
}
