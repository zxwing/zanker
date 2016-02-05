package lib

import (
	"os"
)

func IsFile(path string) bool {
	info, err := os.Stat(s.file)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}
