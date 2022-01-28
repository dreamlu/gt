package fs

import (
	"archive/zip"
	"io"
)

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: 输出的zip文件的名字文件流
// Param 2: 需要添加到zip文件里面的文件文件流
func ZipFiles(newZipFile io.Writer, files []*File) error {

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, f := range files {
		header, err := zip.FileInfoHeader(f)
		if err != nil {
			return err
		}
		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, f); err != nil {
			return err
		}
	}
	return nil
}
