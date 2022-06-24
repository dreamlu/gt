package file

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// test upload
func TestFile_GetUploadFile(t *testing.T) {

	filenameSplit := strings.SplitAfter("test.test.jpg", ".")
	t.Log(filenameSplit)
	fType := filenameSplit[len(filenameSplit)-1]
	t.Log(fType)

	t.Log(filepath.Ext("test.test.jpg"))
}

// 8.9MB->722.3MB
func TestCompressImage(t *testing.T) {
	fileImg := File{
		Path: "../../test/file/呵呵.png",
		//Width:  200,
		//Height: 0,
		//Name: "newName.png",
		NewPath:     "../../test/file/呵呵1.png",
		Quality:     0,
		ContentType: "png",
	}
	err := fileImg.compressImage()
	if err != nil {
		t.Error(err)
	}
}

func TestFileType(t *testing.T) {
	data, _ := ioutil.ReadFile("test.jpg")
	t.Log(ContentType(data))
}

func TestContainsTransparent(t *testing.T) {
	data, _ := ioutil.ReadFile("../../test/file/不透明.png")
	img, _ := png.Decode(bytes.NewReader(data))
	ContainsTransparent(img)
}
