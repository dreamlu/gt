// package gt

package gt

import (
	"github.com/dreamlu/go-tool/tool/result"
	"testing"
)

func TestStr(t *testing.T) {
	//get defalult upload max size
	//r := routers.SetRouter()
	//MaxUploadMemory = r.MaxMultipartMemory
	t.Log("str :", MaxUploadMemory)

	t.Log(result.GetMapData(0,"1").String())
}
