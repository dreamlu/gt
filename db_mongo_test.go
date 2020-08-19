package gt

import (
	"testing"
)

func TestMongoCrud(t *testing.T) {

	var user = User{
		Name: "test",
	}
	cd := NewCrud(
		D("mongo"),
		Table("user"),
		Data(user),
	)
	cd.Create()
}
