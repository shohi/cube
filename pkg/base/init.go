package base

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var (
	DefaultBaseConfigDir string
	DefaultCacheDir      string
	DefaultCertDir       string
	LocalKubeConfigPath  = "~/.kube/config"

	ErrFailedCreateCacheDir = errors.New("failed to create cache dir")
)

func init() {
	var err error
	DefaultBaseConfigDir, err = homedir.Expand("~/.config/cube")
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCacheDir, err))

	}

	DefaultCertDir = filepath.Join(DefaultBaseConfigDir, "cert")
	err = os.MkdirAll(DefaultCertDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCacheDir, err))
	}

	DefaultCacheDir = filepath.Join(DefaultBaseConfigDir, "cache")
	err = os.MkdirAll(DefaultCacheDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCacheDir, err))
	}
}
