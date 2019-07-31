package kube

import clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

// getNextLocalPort get next available local port.
// It checks the cluster whose server is in format `https://kubernetes:xxx`.
func getNextLocalPort(kc *clientcmdapi.Config) (int, error) {
	// TODO

	return 0, nil
}
