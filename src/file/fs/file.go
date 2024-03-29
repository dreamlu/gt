package fs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"time"
)

// File implement file system
// only support file
type File struct {
	name    string
	content *bytes.Buffer
	modTime time.Time
	closed  bool
}

func NewFile() *File {
	return &File{
		content: bytes.NewBuffer(nil),
		modTime: time.Now(),
	}
}

func OpenFile(name string) (*File, error) {
	f, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return &File{
		content: bytes.NewBuffer(f),
		modTime: time.Now(),
	}, nil
}

func (f *File) Write(p []byte) (int, error) {
	if f.closed {
		return 0, errors.New("file closed")
	}

	return f.content.Write(p)
}

func (f *File) Read(p []byte) (int, error) {
	if f.closed {
		return 0, errors.New("file closed")
	}

	return f.content.Read(p)
}

func (f *File) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, errors.New("file closed")
	}

	return f, nil
}

// Close 关闭文件，可以调用多次。
func (f *File) Close() error {
	f.closed = true
	return nil
}

// impl fs.FileInfo

func (f *File) Name() string {
	return f.name
}

func (f *File) SetName(name string) {
	f.name = name
}

func (f *File) Size() int64 {
	return int64(f.content.Len())
}

func (f *File) Mode() fs.FileMode {
	return 0666
}

func (f *File) ModTime() time.Time {
	return f.modTime
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Sys() any {
	return nil
}

func (f *File) Bytes() []byte {
	return f.content.Bytes()
}

func (f *File) WriteTo(w io.Writer) (n int64, err error) {
	if f.content == nil {
		return 0, errors.New("file nil error")
	}
	return f.content.WriteTo(w)
}
