package kube

import (
	"io/ioutil"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

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
