package filechannel

import (
	"testing"
	"os"
)

func init() {
	err := os.RemoveAll("test/chanel")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("test/channel", 0666)
	if err != nil {
		panic(err)
	}
}

func TestFileChannel(t *testing.T) {
	_, err := FileChannel("createTest", "test/channel")
	if err != nil {
		t.Log(err)
	}
}