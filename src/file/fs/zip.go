package fs

import (
	"archive/zip"
	"io"
)

// ZipFiles zip files
func ZipFiles(newZipFile io.Writer, files []*File) (err error) {

	w := zip.NewWriter(newZipFile)
	for _, file := range files {
		var f io.Writer
		f, err = w.Create(file.Name())
		if err != nil {
			return
		}
		_, err = f.Write(file.Bytes())
		if err != nil {
			return
		}
	}
	// Make sure to check the error on Close.
	return w.Close()
}
