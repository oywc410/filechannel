package filechannel

import (
	"os"
	"errors"
	"path/filepath"
	"strconv"
)

type channelFiles struct {
	keyFile *os.File
	dateFile []*os.File
}

func folderCheck(dirpath string) (string, error) {

	if dirpath != "" {
		fileDir, err := os.Stat(dirpath)
		if err != nil {
			if os.IsExist(err) {
				return "", errors.New("Please enter the save folder path")
			}
			return "", err
		}

		if !fileDir.IsDir() {
			return "", errors.New("Please specify a folder")
		}

		return filepath.Abs(dirpath)
	}

	return "", errors.New("Please enter the save folder path")
}

func openFileStart(dirpath string, fileKeyName string) (*channelFiles, error) {
	keyFile, err := os.OpenFile(fileKeyName + ".key_data", os.O_CREATE | os.O_WRONLY | os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	n, err = keyFile.Read(buf)
	if err != nil {
		return nil, err
	}

	fileCount, err := strconv.Atoi(string(buf[:n]))

	if err != nil {
		return nil, err
	}

	return &channelFiles{
		keyFile: keyFile,
	}, nil
}

func openFileClose(cFile *channelFiles) error {
	cFile.keyFile.Close()
	for _, fr := range cFile.dateFile {
		fr.Close()
	}
	return nil
}