package yumlib

import (
	"encoding/xml"
	"fmt"
)

const (
	RepomdFilePath string = "/repodata/repomd.xml"
)

type Location struct {
	Url string `xml:"href,attr"`
}
type Checksum struct {
	Type string `xml:"type,attr"`
	Data string `xml:",chardata"`
}
type RepomdDatabase struct {
	Name         string   `xml:"type,attr"`
	Location     Location `xml:"location"`
	Checksum     Checksum `xml:"checksum"`
	Timestamp    string   `xml:"timestamp"`
	Size         uint64   `xml:"size"`
	OpenSize     uint64   `xml:"open-size"`
	OpenChecksum Checksum `xml:"open-checksum"`
	Version      string   `xml:"version"`
}

func (self *RepomdDatabase) GetChecksum() (digesttype string, data []byte) {
	
	return
}

type Repomd struct {
	Revision string           `xml:"revision"`
	Database []RepomdDatabase `xml:"data"`
}

func NewRepomd(xmldata []byte) (r *Repomd, err error) {
	r = new(Repomd)

	err = xml.Unmarshal(xmldata, r)

	if err == nil {
		return r, nil
	} else {
		return nil, err
	}
}

func (r *Repomd) GetDatabaseRepomd(name string) (repomddb *RepomdDatabase, err error) {
	repomddb = nil
	err = fmt.Errorf("RepomdDatabase name '%s' not found in repomd.xml", name)

	for _, r := range r.Database {
		if r.Name == name {
			repomddb = &r
			err= nil
			break
		}
	}
	return repomddb, err
}
