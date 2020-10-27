package base

import "os"

// FileExists checks if a file exists
// refer, https://golangcode.com/check-if-a-file-exists/
func FileExists(filename string) (exist, isDir bool) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, false
	}

	return true, info.IsDir()
}
