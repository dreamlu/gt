package conf

import (
	"bytes"
	"github.com/dreamlu/gt/tool/util/gos"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ProjectPath return project path
func ProjectPath() (path string) {
	// default linux/mac os
	var (
		sp = "/"
		ss []string
	)
	if runtime.GOOS == "windows" {
		sp = "\\"
	}

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
