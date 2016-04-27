package filechannel

import (

)

type filechannel struct {
	channelKey string
	channePath string
	startFileNo int
	endFileNo int
	startNo int
	endNo int
	channel chan []byte
	channelFiles *channelFiles
}

func FileChannel(channelKey string, channePath string) (*filechannel, error) {

	startNo := 0
	endNo := 0
	startFileNo := 0
	endFileNo := 0

	f := &filechannel{
		channelKey: channelKey,
		channePath: channePath,
		startFileNo: startFileNo,
		endFileNo: endFileNo,
		startNo: startNo,
		endNo: endNo,
		channel: make(chan []byte),
	}

	return f, f.init()
}

func (f *filechannel) init() error {
	var err error
	f.channePath, err = folderCheck(f.channePath)
	if err != nil {
		return err
	}

	f.channelFiles, err = openFileStart(f.channePath, f.channelKey)
	if err != nil {
		return err
	}

	return nil
}

func (f *filechannel) Get() []byte {
	return <- f.channel
}

func (f *filechannel) Set(date []byte) {
	f.channel <- date
}