package conf

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/dreamlu/gt/src/gos"
)

// default linux/mac os
var sp = "/"

func init() {
	if runtime.GOOS == "windows" {
		sp = "\\"
	}
}

// rPath: relative path
// aPath: absolute path
// return config path
func newPath(rPath string) string {
	if workDir, err := os.Getwd(); err == nil {
		wPath := workDir + sp + rPath
		if gos.Exists(wPath) {
			return wPath
		}
	}

	aPath := ProjectPath() + rPath
	if gos.Exists(aPath) {
		return aPath
	}
	return rPath
}

// ProjectPath return project path
func ProjectPath() (path string) {
	var ss []string

	// GOMOD
	// in go source code:
	// // Check for use of modules by 'go env GOMOD',
	// // which reports a go.mod file path if modules are enabled.
	// stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	// gomod := string(bytes.TrimSpace(stdout))
	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	path = string(bytes.TrimSpace(stdout))
	if path != "" {
		ss = strings.Split(path, sp)
		ss = ss[:len(ss)-1]
		path = strings.Join(ss, sp) + sp
		return
	}

	// GOPATH
	fileDir, _ := os.Getwd()
	path = os.Getenv("GOPATH") // < go 1.17 use
	ss = strings.Split(fileDir, path)
	if path != "" {
		ss2 := strings.Split(ss[1], sp)
		path += sp
		for i := 1; i < len(ss2); i++ {
			path += ss2[i] + sp
			if gos.Exists(path) {
				return path
			}
		}
	}
	return
}
