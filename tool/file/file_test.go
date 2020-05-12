package file

import (
	"github.com/dreamlu/gt"
	"os"
	"testing"
)

func TestCompressImage(t *testing.T) {
	fileImg := File{
		Path:    "../../test/file/呵呵.jpg",
		NewPath: "../../test/file/呵呵1.jpg",
		Width:   200,
		Height:  0,
	}
	err := fileImg.CompressImage("jpg")
	if err != nil {
		t.Error(err)
	}
}

// 子目录读取文gt配置文件测试
func TestConfigger(t *testing.T) {
	dir, _ := os.Getwd()
	t.Log(dir)
	mode := gt.Configger().GetString("app.devMode")
	t.Log(mode)
}
