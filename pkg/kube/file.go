package kube

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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

// getLocalPath creates local config path from remote address by convention.
// localPath is `~/.config/cube/cache/$HOST`.
func getLocalPath(remoteAddr string) string {
	filename := extractHost(remoteAddr)

	return filepath.Join(DefaultCacheDir, filename+".yaml")
}

func getLocalKubePath() string {
	p, err := homedir.Expand(LocalKubeConfigPath)
	if err != nil {
		panic(fmt.Sprintf("failed to get kubeconfig locally, err: %v", err))
	}

	return p
}

func getLocalCertAuthPath(remoteAddr string) string {
	filename := extractHost(remoteAddr)

	return filepath.Join(DefaultCertDir, filename+"-ca.crt")
}

func getLocalCertClientPath(remoteAddr string) string {
	filename := extractHost(remoteAddr)

	return filepath.Join(DefaultCertDir, filename+"-client.crt")
}

func getLocalCertClientKeyPath(remoteAddr string) string {
	filename := extractHost(remoteAddr)

	return filepath.Join(DefaultCertDir, filename+"-client.key")
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
