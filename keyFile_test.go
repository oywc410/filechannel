package filechannel

import (
	"testing"
	"path"
	"os"
	"sync"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//文件计数
func TestCountFile(t *testing.T) {

	filePath := path.Join("test", "keyFile.test")

	fs, err := NewKeyFile(filePath)
	if err != nil {
		t.Error(err)
	}

	fileCount, err := fs.GetDataFileCount()
	if err != nil {
		t.Error(err)
	}

	if fileCount != 0 {
		t.Error("err 01")
	}

	_, err = fs.AddDataFileCount()
	if err != nil {
		t.Error(err)
	}

	fs.AddDataFileCount()
	fs.AddDataFileCount()
	fs.AddDataFileCount()

	fileCount, _ = fs.AddDataFileCount()

	if fileCount != 5 {
		t.Errorf("数据文件数:%d", fileCount)
	}
	t.Logf("数据文件数:%d", fileCount)

	err = fs.CloseFile()
	if err != nil {
		t.Error(err)
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Error(err)
	}

}

//并发下的文件计数
func TestConcurrentCount(t *testing.T) {
	var w sync.WaitGroup

	filePath := path.Join("test", "keyFile.test")
	fs, err := NewKeyFile(filePath)
	if err != nil {
		t.Error(err)
	}

	j := 600
	for {
		j--
		if j < 0 {
			break
		}
		w.Add(1)
		go func() {
			i := 10
			for {
				i--
				if i < 0 {
					break
				}
				fs.AddDataFileCount()
			}

			w.Done()
		}()
	}
	w.Wait()

	count, _ := fs.GetDataFileCount()
	if count != 6000 {
		t.Errorf("并发数据文件数:%d", count)
	}
	t.Logf("并发数据文件数:%d", count)

	err = fs.CloseFile()
	if err != nil {
		t.Error(err)
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Error(err)
	}

}
