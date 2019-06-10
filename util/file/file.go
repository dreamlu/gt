// @author  dreamlu
package file

import (
	"github.com/dreamlu/go-tool"
	"github.com/dreamlu/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

// file
type File struct {
}

//获得文件上传路径,内部专用
func (f *File) GetUploadFile(file *multipart.FileHeader, fname string) (filename string) {

	filenameSplit := strings.Split(file.Filename, ".")
	ftype := filenameSplit[len(filenameSplit)-1]
	//防止文件名中多个“.”,获得文件后缀
	filename = "." + ftype
	switch fname {
	case "":                                                      //重命名
		filename = time.Now().Format("20060102150405") + filename //时间戳"2006-01-02 15:04:05"是参考格式,具体数字可变(经测试)
	default: //指定文件名
		//防止文件名中多个“.”,获得文件后缀
		filename = fname + filename
	}
	path := der.GetDevModeConfig("filepath") + filename //文件目录
	_ = f.SaveUploadedFile(file, path)
	switch ftype {
	case "jpeg", "jpg", "png":
		_ = f.CompressImage(ftype, path)
	default:
		//处理其他类型文件
	}

	return path
}

//单文件上传
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

//图片压缩
func (f *File) CompressImage(imagetype, path string) error {
	//图片压缩
	var img image.Image
	ImgFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer ImgFile.Close()

	switch imagetype {
	case "jpeg", "jpg":
		img, err = jpeg.Decode(ImgFile)
		if err != nil {
			return err
		}
	case "png":
		img, err = png.Decode(ImgFile)
		if err != nil {
			return err
		}
	}

	m := resize.Resize(0, 0, img, resize.Lanczos3)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	switch imagetype {
	case "jpeg", "jpg":
		// write new image to file
		_ = jpeg.Encode(out, m, nil)
	case "png":
		_ = png.Encode(out, m) // write new image to file
	}

	return nil
}

func (f *File) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
