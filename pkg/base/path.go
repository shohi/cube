package base

import (
	"fmt"
	"path/filepath"

	"github.com/atrox/homedir"
)

// GetLocalKubePath returns local kubeconfig absolute path.
func GetLocalKubePath() string {
	p, err := homedir.Expand(LocalKubeConfigPath)
	if err != nil {
		panic(fmt.Sprintf("failed to get kubeconfig locally, err: %v", err))
	}

	return p
}

// GenLocalCertAuthPath creates local path for remote cert-auth file
func GenLocalCertAuthPath(remoteAddr string) string {
	filename := ExtractHost(remoteAddr)
	return filepath.Join(DefaultCertDir, filename+"-ca.crt")
}

// GenLocalCertClientPath creates local path for remote client-cert file
func GenLocalCertClientPath(remoteAddr string) string {
	filename := ExtractHost(remoteAddr)

	return filepath.Join(DefaultCertDir, filename+"-client.crt")
}

// GenLocalCertClientKeyPath creates local path for remote client-key file
func GenLocalCertClientKeyPath(remoteAddr string) string {
	filename := ExtractHost(remoteAddr)

	return filepath.Join(DefaultCertDir, filename+"-client.key")
}
