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
	DefaultHistoryPath   string

	LocalKubeConfigPath = "~/.kube/config"

	ErrFailedCreateCacheDir = errors.New("failed to create cache dir")
	ErrFailedCreateCertDir  = errors.New("failed to create cert dir")
	ErrFailedCreateHistory  = errors.New("failed to create history file")
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
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCertDir, err))
	}

	DefaultCacheDir = filepath.Join(DefaultBaseConfigDir, "cache")
	err = os.MkdirAll(DefaultCacheDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCacheDir, err))
	}

	DefaultHistoryPath = filepath.Join(DefaultBaseConfigDir, "history")
	f, err := os.OpenFile(DefaultHistoryPath, os.O_RDONLY|os.O_CREATE, 0666)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateHistory, err))
	}
}
