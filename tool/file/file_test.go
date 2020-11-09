package file

import (
	"bytes"
	"github.com/dreamlu/gt/tool/conf"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// test upload
func TestFile_GetUploadFile(t *testing.T) {
	t.Log(conf.Configger().GetString("app.filepath") + time.Now().Format("20060102") + "/")
}

// 8.9MB->722.3MB
func TestCompressImage(t *testing.T) {
	fileImg := File{
		Path: "../../test/file/呵呵.png",
		//Width:  200,
		//Height: 0,
		NewPath: "../../test/file/呵呵1.png",
		Quality: 0,
	}
	err := fileImg.CompressImage("png")
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

func TestFileType(t *testing.T) {
	data, _ := ioutil.ReadFile("../../test/file/呵呵.jpg")
	contentType := GetFileContentType(data[:512])
	t.Log(contentType)
}

func TestContainsTransparent(t *testing.T) {
	data, _ := ioutil.ReadFile("../../test/file/不透明.png")
	img, _ := png.Decode(bytes.NewReader(data))
	ContainsTransparent(img)
}
