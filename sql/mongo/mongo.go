package mongo

import (
	"context"
	"fmt"
	"github.com/dreamlu/gt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	onceMongoDB sync.Once
	// mongoDB is global
	mongoDB *mongo.Database
)

// db params
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

func newMongoDB() *mongo.Database {

	dbS := &dba{}
	gt.Configger().GetStruct("app.mongo", dbS)
	//url := fmt.Sprintf("mongodb://%s:%s@%s", dbS.User, dbS.Password, dbS.Host)
	url := fmt.Sprintf("mongodb://%s", dbS.Host)
	ctx, _ := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		gt.Logger().Error("[mongodb连接错误]:", err)
		gt.Logger().Error("[mongodb开始尝试重连中]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			client, err = mongo.Connect(ctx, options.Client().ApplyURI(url))
			//defer DB.Close()
			if err != nil {
				gt.Logger().Error("[mongodb连接错误]:", err)
				continue
			}
			gt.Logger().Info("[mongodb重连成功]")
			break
		}
	}

	return client.Database(gt.Configger().GetString("app.mongo.name"))
}

// single db
func MongoDB() *mongo.Database {

	onceMongoDB.Do(func() {
		mongoDB = newMongoDB()
	})
	return mongoDB
}
