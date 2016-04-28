package filechannel

import (
	"os"
	"testing"
	"path/filepath"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := os.RemoveAll(filepath.Join("test", "file_test"))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join("test", "file_test"), 0666)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filepath.Join("test", "file_test.txt"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
}

func TestFolderCheck(t *testing.T) {
	_, err := folderCheck(filepath.Join("test", "aaaa"))
	if err == nil {
		t.Error("err 01")
	}

	t.Log(err)

	dirname, err := folderCheck(filepath.Join("test", "file_test"))
	if err != nil {
		t.Error("err 02")
	}

	t.Log(dirname)

	_, err = folderCheck(filepath.Join("test", "file_test.txt"))
	if err == nil {
		t.Error("err 03")
	}

	t.Log(err)
}

