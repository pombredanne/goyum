package yumlib

import (
	"io/ioutil"
	"testing"
)

const (
	TestDatabaseFile     string = "testdata/test.sqlite.bz2"
	TestDatabaseOpenSize uint64 = 19504128
	TestDatabaseRows     int    = 6019
)

func TestPrimaryDatabase(t *testing.T) {
	data, err := ioutil.ReadFile(TestDatabaseFile)
	if err != nil {
		t.Fatal(err)
	}

	db, err := NewDatabase(data, "testdb", TestDatabaseOpenSize)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM packages")
	if err != nil {
		t.Fatal(err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}

	if counter != TestDatabaseRows {
		t.Errorf("Query result not matched expected count = %d actual = %d\n",
			TestDatabaseRows, counter)
	}
}
