package file

import (
	"bytes"
	"errors"
	"github.com/dreamlu/gt/serv/conf"
	"github.com/dreamlu/gt/serv/snowflake"
	"github.com/dreamlu/gt/tool/cons"
	"github.com/dreamlu/gt/tool/gos"
	"github.com/dreamlu/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// example
// use gin upload file
// func Upload(u *gin.Context) {
//
//	name := u.PostForm("name") // set file name
//	file, err := u.FormFile("file")
//	if err != nil {
//		u.JSON(http.StatusOK, someErr)
//	}
//	path := File{Name: name}.GetUploadFile(file)
//	u.JSON(http.StatusOK, path)
//}

const (
	JPEG = "jpeg" // jpeg/jpg
	PNG  = "png"  // png
)

// File upload
type File struct {
	File *multipart.FileHeader
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
	// is img compress
	// default false, no compress
	IsComp  bool
	Quality int // default 80, 1-100

	// content type
	ContentType string
}

// NewFile
// file sugar
func NewFile(file *multipart.FileHeader, Name string) *File {
	return &File{File: file, Name: Name}
}

// Upload file
func (f *File) Upload() (err error) {
	fileExt := filepath.Ext(f.File.Filename)
	switch f.Name {
	case "":
		snowflakeID, err := snowflake.NewID(1)
		if err != nil {
			return err
		}
		f.Name = snowflakeID.String() + fileExt
	default:
		f.Name += fileExt
	}
	err = f.Save()
	if err != nil {
		return err
	}

	// whatever
	if f.IsComp {
		go f.Compress()
	}
	return nil
}

// Save file
func (f *File) Save() (err error) {

	if f.Format == "" {
		f.Format = "20060102"
	}
	f.Path = conf.Get[string](cons.ConfFile) + time.Now().Format(f.Format) + "/"
	if err = gos.Mkdir(f.Path); err != nil {
		return
	}

	// File Path
	f.Path += f.Name
	src, err := f.File.Open()
	if err != nil {
		return
	}
	defer src.Close()

	out, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return
}

// ImageConfig return Upload Image Config width/height
func (f *File) ImageConfig() (*image.Config, error) {

	if f.IsImg() {
		f, _ := os.Open(f.Path)
		defer f.Close()
		c, _, err := image.DecodeConfig(f)
		return &c, err
	}
	return nil, errors.New("not Image type")
}

// Compress Image
func (f *File) Compress() {
	if f.IsImg() {
		_ = f.compressImage()
	}
}

func (f *File) IsImg() bool {
	if f.ContentType == "" {
		data, _ := ioutil.ReadFile(f.Path)
		f.ContentType = GetImageType(data)
	}
	if strings.Contains(f.ContentType, PNG) || strings.Contains(f.ContentType, JPEG) {
		return true
	}
	return false
}

// compressImage image compress
func (f *File) compressImage() error {
	var img image.Image
	imgFile, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}

	switch f.ContentType {
	case JPEG:
		img, err = jpeg.Decode(bytes.NewReader(imgFile))
	case PNG:
		img, err = png.Decode(bytes.NewReader(imgFile))
	default:
		return errors.New("[gt] not support img type:" + f.ContentType)
	}
	if err != nil {
		return err
	}

	if f.NewPath != "" {
		f.Path = f.NewPath
	}

	m := resize.Resize(uint(f.Width), uint(f.Height), img, resize.Lanczos3)

	out, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer out.Close()

	switch f.ContentType {
	case JPEG:
		// write new image to file
		var q *jpeg.Options
		if f.Quality > 0 {
			q = &jpeg.Options{Quality: f.Quality}
		}
		_ = jpeg.Encode(out, m, q)
	case PNG:
		if ContainsTransparent(m) {
			_ = png.Encode(out, m) // write new image to file
		} else {
			_ = PngToJpeg(m, out, f.Quality)
		}
	}

	return nil
}

// GetImageType jpeg,png
func GetImageType(buffer []byte) string {
	contentType := GetFileContentType(buffer)

	switch contentType {
	case "image/jpeg":
		return JPEG
	case "image/png":
		return PNG
	default:
		return ""
	}
}

// GetFileContentType must a file
// file byte data[:512]
// image type: "image/jpeg","image/png"
func GetFileContentType(buffer []byte) string {
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:512])
	return contentType
}
