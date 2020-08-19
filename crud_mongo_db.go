package gt

import (
	"context"
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

func newMongoDB() *mongo.Database {
	ctx, _ := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+Configger().GetString("app.mongo.host")))
	if err != nil {
		Logger().Error("[mongodb连接错误]:", err)
		Logger().Error("[mongodb开始尝试重连中]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+Configger().GetString("app.mongo.host")))
			//defer DB.Close()
			if err != nil {
				Logger().Error("[mongodb连接错误]:", err)
				continue
			}
			Logger().Info("[mongodb重连成功]")
			break
		}
	}

	return client.Database(Configger().GetString("app.mongo.name"))
}

// single db
func MongoDB() *mongo.Database {

	onceMongoDB.Do(func() {
		mongoDB = newMongoDB()
	})
	return mongoDB
}
