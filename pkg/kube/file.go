package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultCacheDir = "~/.config/cube/cache"
)

func init() {
	err := os.MkdirAll(DefaultCacheDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("failed to create cache dir, err: %v", err))
	}
}

// extractHost extracts host info from remoteAddr which is in the format `user@host`
func extractHost(remoteAddr string) string {
	tokens := strings.Split(remoteAddr, "@")

	return tokens[len(tokens)-1]
}

// getLocalPath creates local config path from remote address by convention.
// localPath is `~/.config/cube/cache/$HOST`.
func getLocalPath(remoteAddr string) string {
	filename := extractHost(remoteAddr)

	return filepath.Join(DefaultCacheDir, filename)
}
