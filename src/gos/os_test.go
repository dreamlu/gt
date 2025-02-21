package gos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	t.Log(Exists("/"), Exists("/test"))
	f, _ := os.Open("/")
	t.Log(f.Readdirnames(0))
	t.Log(filepath.Abs(""))
}

func TestMkFile(t *testing.T) {
	t.Log(MkFile("./tmp/test.txt"))
}
