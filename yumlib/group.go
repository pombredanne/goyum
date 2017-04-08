package yumlib

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io"
)

type GroupPackageData struct {
	Requires string `xml:"requires,attr"`
	Type     string `xml:"type,attr"`
	Name     string `xml:",chardata"`
}

type GroupPackageList struct {
	PackageList []GroupPackageData `xml:"packagereq"`
}

type GroupName struct {
	Lang string `xml:"lang,attr"`
	Name string `xml:",chardata"`
}

type GroupDescription struct {
	Lang        string `xml:"lang,attr"`
	Description string `xml:",chardata"`
}

type GroupData struct {
	Id              string           `xml:"id"`
	Names           []GroupName      `xml:"name"`
	Default         bool             `xml:"default"`
	Description     GroupDescription `xml:"description"`
	PackageListData GroupPackageList `xml:"packagelist"`
}

type GroupListData struct {
	Groups []GroupData `xml:"group"`
}

func NewGroupListData(gzdata []byte) (grouplist *GroupListData, err error) {
	reader, err := gzip.NewReader(bytes.NewReader(gzdata))
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	var xmldata []byte

	for {
		buffer := make([]byte, 512)

		n, err := reader.Read(buffer)
		if err == io.EOF {
			if n > 0 {
				xmldata = append(xmldata, buffer[:n]...)
			}
			break
		}
		if err != nil {
			return nil, err
		}

		xmldata = append(xmldata, buffer[:n]...)
	}

	grouplist = new(GroupListData)

	err = xml.Unmarshal(xmldata, grouplist)
	if err != nil {
		return nil, err
	}

	return grouplist, nil

}
