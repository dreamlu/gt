package tag

import (
	"testing"
)

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
		OtherX string `json:"other_x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name"`
		UserDeX
	}
	// test tag
	t.Log(GetTags(UserMore{}))
}
