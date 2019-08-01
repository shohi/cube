package scp

import "os"

// hasContent checks whether file exists and is nonempty.
func hasContent(filename string) bool {
	fi, err := os.Stat(filename)
	if err == nil && fi.Size() > 0 {
		return true
	}

	return false
}
