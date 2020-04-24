package file

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

// example
// 单文件上传
// use gin upload file
//func UpoadFile(u *gin.Context) {
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
		switch strings.ToLower(fType) {
		case "jpeg", "jpg", "png":
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
	ImgFile, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer ImgFile.Close()

	switch imageType {
	case "jpeg", "jpg":
		img, err = jpeg.Decode(ImgFile)
	case "png":
		img, err = png.Decode(ImgFile)
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
	case "jpeg", "jpg":
		// write new image to file
		_ = jpeg.Encode(out, m, nil)
	case "png":
		_ = png.Encode(out, m) // write new image to file
	}

	return nil
}

// save file
func (f *File) SaveUploadedFile(file *multipart.FileHeader, filename string) (path string, err error) {

	filepath := gt.Configger().GetString("app.filepath")
	if !Exists(filepath) {
		err = os.Mkdir(filepath, os.ModePerm)
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
