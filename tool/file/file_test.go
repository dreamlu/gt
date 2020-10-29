package file

import (
	"github.com/dreamlu/gt/tool/conf"
	"os"
	"testing"
	"time"
)

// test upload
func TestFile_GetUploadFile(t *testing.T) {
	t.Log(conf.Configger().GetString("app.filepath") + time.Now().Format("20060102") + "/")
}

func TestCompressImage(t *testing.T) {
	fileImg := File{
		Path:   "../../test/file/呵呵.jpg",
		Width:  200,
		Height: 0,
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
	mode := conf.Configger().GetString("app.devMode")
	t.Log(mode)
}
