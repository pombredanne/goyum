package yumlib

import (
	"net/http"
	"io/ioutil"
)

type Downloader interface {
	Get(location string) (data []byte, err error)
}

type HttpDownloader struct {
}

func NewHttpDownloader() (downloader Downloader) {
	downloader = new(HttpDownloader)

	return downloader
}

func (self *HttpDownloader) Get(location string) (data []byte, err error) {
	res, err := http.Get(location)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	
	return data, err
}

type LocalFileDownloader struct {
}

func NewLocalFileDownloader() (downloader Downloader) {
	downloader = new(LocalFileDownloader)
	
	return downloader
}

func (self *LocalFileDownloader) Get(location string) (data []byte, err error) {
	data, err = ioutil.ReadFile(location)
	
	return
}