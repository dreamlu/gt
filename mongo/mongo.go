package mongo

import (
	"context"
	"fmt"
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/log"
	"github.com/dreamlu/gt/src/cons"
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
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Name        string `yaml:"name"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
	MaxOpenConn int    `yaml:"maxOpenConn"`
	// db log mode
	Log bool `yaml:"log"`
}

func newMongoDB() *mongo.Database {

	dbS := &dba{}
	conf.UnmarshalField(cons.ConfMongo, dbS)
	//url := fmt.Sprintf("mongodb://%s:%s@%s", dbS.User, dbS.Password, dbS.Host)
	url := fmt.Sprintf("mongodb://%s", dbS.Host)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Error("[mongodb connect error]:", err)
		log.Error("[mongodb try connect again]: try it every 5s...")
		// try to reconnect
		for {
			// go is so fast
			// try it every 5s
			time.Sleep(5 * time.Second)
			client, err = mongo.Connect(ctx, options.Client().ApplyURI(url))
			//defer DB.Close()
			if err != nil {
				log.Error("[mongodb connect error]:", err)
				continue
			}
			log.Info("[mongodb connect success]")
			break
		}
	}

	return client.Database(dbS.Name)
}

// MongoDB single db
func MongoDB() *mongo.Database {

	onceMongoDB.Do(func() {
		mongoDB = newMongoDB()
	})
	return mongoDB
}
