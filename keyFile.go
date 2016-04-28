package filechannel

import (
	"os"
	"encoding/binary"
	"sync"
)

type keyFile struct {
	filePath string
	fr *os.File
	fileCountLock sync.RWMutex
}

func NewKeyFile(filePath string) (*keyFile, error) {

	fr, err := os.OpenFile(filePath, os.O_CREATE | os.O_WRONLY | os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	return &keyFile{
		filePath: filePath,
		fr: fr,
	}, nil
}

func (f *keyFile)GetDataFileCount() (uint64, error) {
	f.fileCountLock.RLock()
	defer f.fileCountLock.RUnlock()
	f.fr.Seek(0)

	buf := make([]byte, 8)
	n, err := f.fr.Read(buf)
	if err != nil {
		return 0, nil
	}

	return binary.BigEndian.Uint64(buf[:n]), nil
}

func (f *keyFile)AddDataFileCount() error {
	f.fileCountLock.Lock()
	defer f.fileCountLock.Unlock()
	count, err := f.GetDataFileCount()
	count++
	if err != nil {
		return 0, err
	}

	f.fr.Seek(0)
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, count)
	_, err = f.fr.Write(buf)
	return err
}


func (f *keyFile)CloseFile() error {
	return f.fr.Close()
}