package kube

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var DefaultCacheDir string
var ErrFailedCreateCacheDir = errors.New("failed to create cache dir")

func init() {
	var err error
	DefaultCacheDir, err = homedir.Expand("~/.config/cube/cache")
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCacheDir, err))

	}

	err = os.MkdirAll(DefaultCacheDir, os.ModePerm)
	log.Printf("init cache dir")
	if err != nil {
		panic(fmt.Sprintf("%v, cause: %v", ErrFailedCreateCacheDir, err))
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

	return filepath.Join(DefaultCacheDir, filename+".yaml")
}

func getLocalKubePath() string {
	p, err := homedir.Expand("~/.kube/config")
	if err != nil {
		panic(fmt.Sprintf("failed to get kubeconfig locally, err: %v", err))
	}

	return p
}

// load reads kubeconfig from file
func load(configPath string) (*clientcmdapi.Config, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	kc, err := clientcmd.Load(content)
	if err != nil {
		return nil, err
	}

	return kc, nil
}
