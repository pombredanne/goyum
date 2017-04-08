package yumlib

import (
	"os"
	"path/filepath"
)

func DefaultTempDirectoryPath() (path string) {
	path = filepath.Join(os.TempDir(), "yum")
	return	
}

func DefaultCacheDirectory() (path string) {
	path = filepath.Join(DefaultTempDirectoryPath(), "cache")
	return
}

func DefaultDBFilePermission() os.FileMode {
	return 0666
}
