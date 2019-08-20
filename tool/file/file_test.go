package file

import (
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
