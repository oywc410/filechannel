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

func (f *keyFile)getDate(index int64) (uint64, error) {
	f.readFileLock.Lock()
	defer f.readFileLock.Unlock()

	f.readFile.Seek(index, 0)
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

func (f *keyFile)setDate(index int64, n uint64) error {
	f.writeFileLock.Lock()
	defer f.writeFileLock.Unlock()

	buf := make([]byte, 8)
	f.writeFile.Seek(index, 0)
	binary.BigEndian.PutUint64(buf, n)
	_, err := f.writeFile.WriteAt(buf, 0)
	if err != nil {
		return err
	}
	//return f.writeFile.Sync()
	return nil
}

func (f *keyFile)nextDate(index int64) (uint64, error) {
	f.writeFileLock.Lock()
	defer f.writeFileLock.Unlock()

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
	if err != nil {
		return 0, err
	}
	//return count, f.writeFile.Sync()
	return count, nil
}

func (f *keyFile)flushAllDate() error {
	return f.writeFile.Sync()
}

func (f *keyFile)GetDataFileCount() (uint64, error) {
	return f.getDate(0)
}

func (f *keyFile)AddDataFileCount() (uint64, error) {
	return f.nextDate(0)
}

func (f *keyFile)GetStartFileIndex() (uint64, error) {
	return f.getDate(1)
}

func (f *keyFile)AddStartFileIndex() (uint64, error) {
	return f.nextDate(1)
}


func (f *keyFile)CloseFile() error {

	f.writeFileLock.Lock()

	defer f.writeFileLock.Unlock()

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