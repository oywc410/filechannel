package filechannel

import (
	"os"
	"encoding/binary"
	"sync"
	"io"
)

type keyFile struct {
	filePath string
	readFile *os.File
	writeFile *os.File
	fileCountLock sync.RWMutex
	readFileLock sync.Mutex
	writeFileLock sync.Mutex
}

func NewKeyFile(filePath string) (*keyFile, error) {

	writeFile, err := os.OpenFile(filePath, os.O_CREATE | os.O_RDWR , 0644)
	if err != nil {
		return nil, err
	}

	readFile, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	return &keyFile{
		filePath: filePath,
		readFile: readFile,
		writeFile: writeFile,
	}, nil
}

func (f *keyFile)GetDataFileCount() (uint64, error) {
	f.fileCountLock.RLock()
	f.readFileLock.Lock()
	defer func() {
		f.readFileLock.Unlock()
		f.fileCountLock.RUnlock()
	}()
	f.readFile.Seek(0, 0)

	buf := make([]byte, 8)
	n, err := f.readFile.Read(buf)
	if err != nil {
		if err == io.EOF {
			return 0, nil
		}
		return 0, err
	}

	return binary.BigEndian.Uint64(buf[:n]), nil
}

func (f *keyFile)AddDataFileCount() (uint64, error) {
	f.fileCountLock.Lock()
	f.writeFileLock.Lock()
	defer func() {
		f.writeFileLock.Unlock()
		f.fileCountLock.Unlock()
	}()

	f.writeFile.Seek(0, 0)
	buf := make([]byte, 8)
	n, err := f.writeFile.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	var count uint64
	if err != io.EOF {
		count = binary.BigEndian.Uint64(buf[:n])
	}
	count++

	binary.BigEndian.PutUint64(buf, count)
	f.writeFile.Seek(0, 0)
	_, err = f.writeFile.WriteAt(buf, 0)

	return count, err
}

func (f *keyFile)CloseFile() error {
	f.fileCountLock.Lock()
	f.writeFileLock.Lock()

	defer func() {
		f.writeFileLock.Unlock()
		f.fileCountLock.Unlock()
	}()

	var err error
	err = f.readFile.Close()
	if err != nil {
		return err
	}

	err = f.writeFile.Close()
	if err != nil {
		return err
	}

	return nil
}