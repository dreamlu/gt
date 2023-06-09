package tag

import (
	"testing"
)

type User struct {
	Name string `json:"name" gt:"field:user.name;like;soft_del"`
}

type UserD struct {
	User
	Other string `json:"other" gt:"ignore;soft_del"`
}

func (User) TableName() string {
	return "users"
}

func TestParseGt(t *testing.T) {
	// test tag
	t.Log(ParseGt(UserD{}))
}

func TestGetTags(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}
	type UserDe struct {
		User
		Other string `json:"other"`
	}

	type UserDeX struct {
		a []string `gt:"ignore"`
		UserDe
		OtherX string `json:"other_x" gt:"field:other.x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name" gorm:"column:shop.name"`
		UserDeX
	}
	// test tag
	t.Log(GetTags(UserMore{}))
}
