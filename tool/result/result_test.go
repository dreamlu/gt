package result

import (
	"testing"
)

func TestMapData_Add(t *testing.T) {
	t.Log(MapCreate.Add("id", 2).Add("test", "3"))
	t.Log(MapUpdate.Add("id", 7).Add("test", "3eyw"))

	pager := GetInfoPager{
		GetInfo: &GetInfo{
			MapData: MapCreate,
			Data:    "test",
		},
		Pager: Pager{
			ClientPage: 2,
			EveryPage:  3,
			TotalNum:   5,
		},
	}
	// pager
	t.Log(pager.Add("id", 1).Add("test", 2))
}
