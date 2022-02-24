package mock

import (
	"github.com/bxcodec/faker/v3"
	"github.com/dreamlu/gt/tool/type/json"
	"github.com/dreamlu/gt/tool/type/time"
	"testing"
)

type User struct {
	ID         uint64     `json:"id"`
	Name       string     `json:"name"`
	BirthDate  time.CDate `json:"birth_date" gorm:"type:date"` // data
	CreateTime time.CTime `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"create_time"`
	Info       json.CJSON `json:"info"`
	Str        []string   `json:"str"`
}

func TestMock(t *testing.T) {
	faker.SetRandomMapAndSliceSize(11)
	var user User
	Mock(&user)
	t.Log(user)
	t.Log(user.Info)
}
