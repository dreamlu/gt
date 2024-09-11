package gos

import (
	"fmt"
	"github.com/dreamlu/gt/src/type/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Exists file/dir exit
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir is dir
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile is file
func IsFile(path string) bool {
	return !IsDir(path)
}

// Mkdir create dir if not exist
func Mkdir(dir string) error {
	if !Exists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil { //os.ModePerm
			return err
		}
	}
	return nil
}

// CopyFile use io.Copy copy file
func CopyFile(src, des string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	// copy all permissions of the source file
	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

// CopyDir copy dir
func CopyDir(srcPath, dstPath string) error {

	if srcPath == dstPath {
		return errors.New(fmt.Sprintf("%s can not the same as %s", srcPath, dstPath))
	}
	if !IsDir(srcPath) || !IsDir(dstPath) {
		return errors.New(fmt.Sprintf("%s or %s is not directory", srcPath, dstPath))
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if path == srcPath {
			return nil
		}
		destNewPath := strings.Replace(path, srcPath, dstPath, -1)
		if !f.IsDir() {
			_, err = CopyFile(path, destNewPath)
		} else {
			return Mkdir(destNewPath)
		}
		return nil
	})

	return err
}
