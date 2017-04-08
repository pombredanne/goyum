package yumlib

import (
	"bytes"
	"compress/bzip2"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewDatabase(
	bz2data []byte, filename string, openSize uint64) (db *sql.DB, err error) {

	decompressed := make([]byte, openSize)

	reader := bzip2.NewReader(bytes.NewReader(bz2data))

	size := 0
	for {
		n, err := reader.Read(decompressed[size:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		size += n

		if uint64(size) >= openSize {
			break
		}
	}

	if uint64(size) != openSize {
		return nil, fmt.Errorf("Decompressed size is invalid actual %d != %d", openSize, size)
	}

	err = os.MkdirAll(DefaultCacheDirectory(), 0755)
	if err != nil {
		return nil, err
	}

	destination := filepath.Join(DefaultCacheDirectory(), filename)

	err = ioutil.WriteFile(destination, decompressed, 0644)
	if err != nil {
		return nil, err
	}

	db, err = sql.Open("sqlite3", destination)

	return db, err
}
