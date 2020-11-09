package file

import (
	"bytes"
	"errors"
	"github.com/dreamlu/gt/tool/conf"
	"github.com/dreamlu/gt/tool/file/file_func"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// example
// 单文件上传
// use gin upload file
//file_func UpoadFile(u *gin.Context) {
//
//	fname := u.PostForm("fname") //指定文件名
//	file, err := u.FormFile("file")
//	if err != nil {
//		u.JSON(http.StatusOK, lib.MapData{Status: lib.CodeFile, Msg: err.Error()})
//	}
//	path := File{}.GetUploadFile(file, fname)
//	u.JSON(http.StatusOK, map[string]interface{}{lib.Status: lib.CodeFile, lib.Msg: lib.MsgFile, "path": path})
//}

// file
type File struct {
	// file name
	Name string
	// path
	Path    string
	NewPath string
	// img attributes
	Width  int
	Height int

	// format 2006-01-02 15:04:05
	Format string

	IsComp  int8 // is img compress
	Quality int  // default 80, 1-100
}

// 获得文件上传路径
func (f *File) GetUploadFile(file *multipart.FileHeader) (filename string, err error) {

	filenameSplit := strings.Split(file.Filename, ".")
	fType := filenameSplit[len(filenameSplit)-1]
	//防止文件名中多个“.”,获得文件后缀
	filename = "." + fType
	switch f.Name {
	case "": //重命名
		snowflakeID, err := id.NewID(1)
		if err != nil {
			return "", err
		}
		filename = snowflakeID.String() + filename //时间戳"2006-01-02 15:04:05"是参考格式,具体数字可变(经测试)
	default: //指定文件名
		//防止文件名中多个“.”,获得文件后缀
		filename = f.Name + filename
	}
	path, err := f.SaveUploadedFile(file, filename)
	if err != nil {
		return "", err
	}

	// whatever
	go func() {
		if f.IsComp == 0 {
			return
		}
		data, _ := ioutil.ReadFile(path)
		fType = GetImageType(data[:512])
		switch fType {
		case "jpeg", "png":
			f.Path = path
			_ = f.CompressImage(fType)
		default:
			//处理其他类型文件
		}
	}()
	return path, nil
}

// 图片压缩
func (f *File) CompressImage(imageType string) error {
	//图片压缩
	var img image.Image
	//imgFile, err := os.Open(f.Path), jpeg.Decode(imgFile)
	imgFile, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}
	//defer ImgFile.Close()

	switch imageType {
	case "jpeg":
		img, err = jpeg.Decode(bytes.NewReader(imgFile))
	case "png":
		img, err = png.Decode(bytes.NewReader(imgFile))
	default:
		return errors.New("[gt] not support img type:" + imageType)
	}
	if err != nil {
		return err
	}

	if f.NewPath == "" {
		f.NewPath = f.Path
	}

	m := resize.Resize(uint(f.Width), uint(f.Height), img, resize.Lanczos3)

	out, err := os.Create(f.NewPath)
	if err != nil {
		return err
	}
	defer out.Close()

	switch imageType {
	case "jpeg":
		// write new image to file
		var q *jpeg.Options
		if f.Quality > 0 {
			q = &jpeg.Options{Quality: f.Quality}
		}
		_ = jpeg.Encode(out, m, q)
	case "png":
		if ContainsTransparent(m) {
			_ = png.Encode(out, m) // write new image to file
		} else {
			_ = PngToJpeg(m, out, f.Quality)
		}
	}

	return nil
}

// save file
func (f *File) SaveUploadedFile(file *multipart.FileHeader, filename string) (path string, err error) {

	if f.Format == "" {
		f.Format = "20060102"
	}
	filepath := conf.Configger().GetString("app.filepath") + time.Now().Format(f.Format) + "/"
	if !file_func.Exists(filepath) {
		err = os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			return
		}
	}

	path = filepath + filename //文件目录
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return
}

// jpeg,png
func GetImageType(buffer []byte) string {
	contentType := GetFileContentType(buffer)

	switch contentType {
	case "image/jpeg":
		return "jpeg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

// file byte data[:512]
// image type: "image/jpeg","image/png"
func GetFileContentType(buffer []byte) string {
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType
}
