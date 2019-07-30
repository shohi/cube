package kube

import (
	"io/ioutil"

	"github.com/shohi/cube/pkg/config"
	"github.com/shohi/cube/pkg/scp"
	"gopkg.in/yaml.v2"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	DefaultKubeConfigPath = "~/.kube/config"
)

// Fuse merges kubeconfig of remote cluster into local
func Fuse(conf config.Config) error {
	localpath := getLocalPath(conf.RemoteAddr)

	// TODO: check whether the file is empty
	err := scp.TransferFile(scp.TransferConfig{
		LocalPath:  localpath,
		RemoteAddr: conf.RemoteAddr,
		RemotePath: DefaultKubeConfigPath,
	})

	if err != nil {
		return err
	}

	//
	_ = clientcmdapi.Config{}

	// TODO: handle duplicated
	// 2. Merge file
	kc, err := newKubeConfig(DefaultKubeConfigPath)
	if err != nil {
		return err
	}

	err = kc.merge(localpath)
	if err != nil {
		return err
	}

	// 3. Print SSH forwarding setting
	// TODO: implement dry-run feature

	// 4. dump file

	return nil
}

// TODO
func printPortForwarding() {

	/*
		if server := os.Getenv("KUBERNETES_MASTER"); len(server) > 0 {
			return server
		}
	*/

}

type KubeConfig struct {
	clientcmdapi.Config
}

func newKubeConfig(conf string) (*KubeConfig, error) {
	content, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}

	var kc clientcmdapi.Config
	err = yaml.Unmarshal(content, &kc)

	if err != nil {
		return nil, err
	}

	return &KubeConfig{Config: kc}, nil
}

func (k *KubeConfig) merge(conf string) error {

	return nil
}
