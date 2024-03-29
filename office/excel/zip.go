package excel

import (
	"github.com/dreamlu/gt/src/file/fs"
	"io"
)

func exportZip[T comparable](dst io.Writer, excels []*Excel[T]) (err error) {

	var fss []*fs.File
	for _, excel := range excels {
		ts := fs.NewFile()
		_, err = excel.WriteTo(ts)
		if err != nil {
			return
		}
		ts.SetName(excel.FileName)
		fss = append(fss, ts)
	}
	return fs.ZipFiles(dst, fss)
}
