package fs

import (
	"fmt"
	"os"
	"testing"
)

func TestZipFiles(t *testing.T) {

	var (
		zip, _ = os.Create("test.zip") // or use io.Writer: NewFile(),it will return io.Writer
		fs     []*File
	)

	for i := 1; i <= 2; i++ {
		ts := NewFile()
		ft, err := OpenFile(fmt.Sprintf("test%d.txt", i))
		if err != nil {
			t.Error(err)
			return
		}
		_, err = ft.WriteTo(ts)
		if err != nil {
			t.Error(err)
			return
		}
		ts.SetName(fmt.Sprintf("test%d_new.txt", i))
		fs = append(fs, ts)
	}

	//zip.SetName("test.zip")
	t.Log(ZipFiles(zip, fs))
}
