package kube

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/shohi/cube/pkg/base"
	"github.com/shohi/cube/pkg/scp"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// DefaultKubeConfigPath is default kubeconfig path on remote host.
	DefaultKubeConfigPath = "~/.kube/config"
)

var (
	ErrRemoteInvalidUser = errors.New("cube: user in remote kubeconfig is invalid")
	ErrRemoteInvalidCert = errors.New("cube: cert in remote kubeconfig is invalid")
	ErrConfigInvalid     = errors.New("cube: remote kubeconfig must have only one cluster")
)

// Downloader download kubernetes config for remote cluster.
// also download cert files if necessary.
type Downloader struct {
	remoteAddr string
	hostIP     string

	kc *clientcmdapi.Config

	clusterName string
	ck          ClusterKeyInfo // key info of cluster who matches given remote addr
}

// DownloadResult represents the download status
type DownloadResult struct {
	ClusterName string               // matched cluster name
	Kc          *clientcmdapi.Config // remote kubectl config

}

// NewDownloader create a new remote config downloader.
func NewDownloader(remoteAddr string) *Downloader {
	hostIP := base.ExtractHost(remoteAddr)
	return &Downloader{
		remoteAddr: remoteAddr,
		hostIP:     hostIP,
	}
}

var emptyDownloadResult = DownloadResult{}

// Download fetches config and cert files.
func (d *Downloader) Download() (DownloadResult, error) {
	if err := d.downloadK8sConfig(); err != nil {
		return emptyDownloadResult, err
	}

	if err := d.checkCertFiles(); err != nil {
		return emptyDownloadResult, err
	}

	result := DownloadResult{
		Kc:          d.kc,
		ClusterName: d.clusterName,
	}

	return result, nil
}

func (d *Downloader) downloadK8sConfig() error {
	p := LocalCachePath(d.remoteAddr)

	// TODO: check whether the file is empty
	err := scp.TransferFile(scp.TransferConfig{
		LocalPath:  p,
		RemoteAddr: d.remoteAddr,
		RemotePath: DefaultKubeConfigPath,
	})
	if err != nil {
		return err
	}

	kc, err := Load(p)
	if err != nil {
		return err
	}

	if err = d.filterCluster(kc); err != nil {
		return err
	}

	d.ck = getClusterKeyInfo(d.kc, d.clusterName)

	return nil
}

// filterCluster gets the matched cluster if multiple clusters exit in
// the Config. The matched cluster is the one whose cluster server hostip
// is equal to the provided hostip. If both http and https exits, prefer
// http one for performance.
func (d *Downloader) filterCluster(kc *clientcmdapi.Config) error {
	if len(kc.Clusters) == 0 {
		return ErrConfigInvalid
	}

	d.kc = kc

	// return immediately if only one cluster is available
	if len(kc.Clusters) == 1 {
		for k := range kc.Clusters {
			d.clusterName = k
			break
		}
		return nil
	}

	var tlsCluster string
	var cluster string
	for k, v := range kc.Clusters {
		if !strings.Contains(v.Server, d.hostIP) {
			continue
		}

		switch {
		case strings.HasPrefix(v.Server, "http://"):
			cluster = k
		case strings.HasPrefix(v.Server, "https://"):
			tlsCluster = k
		}
	}

	// prefer http over https for performance concern
	switch {
	case cluster != "":
		d.clusterName = cluster
		return nil
	case tlsCluster != "":
		d.clusterName = tlsCluster
		return nil
	default:
		return ErrConfigInvalid
	}
}

func (d *Downloader) checkCertFiles() error {
	if d.kc == nil {
		return ErrConfigInvalid
	}

	// TODO: use logrus
	// log.Printf("=====> cluster: [%v]\n", d.cluster)

	// check if token set. If token is set, no cert file/data is needed.
	if len(d.ck.User.Token) > 0 {
		return nil
	}

	// k8s > 1.7
	if len(d.ck.Cluster.CertificateAuthorityData) > 0 {
		if len(d.ck.User.ClientCertificateData) == 0 ||
			len(d.ck.User.ClientKeyData) == 0 {
			return ErrRemoteInvalidUser
		}
		return nil
	}

	// k8s <= 1.7
	if len(d.ck.Cluster.CertificateAuthority) > 0 {
		// Download cert files
		return d.downloadCertFiles()
	}

	return ErrRemoteInvalidCert
}

func (d *Downloader) downloadCertFiles() error {
	if len(d.ck.Cluster.CertificateAuthority) == 0 {
		return nil
	}

	// download auth cert and also update corresponding info
	localAuthPath := base.GenLocalCertAuthPath(d.remoteAddr)
	err := scp.TransferFile(scp.TransferConfig{
		LocalPath:  localAuthPath,
		RemoteAddr: d.remoteAddr,
		RemotePath: d.ck.Cluster.CertificateAuthority,
	})

	if err != nil {
		return err
	}
	d.ck.Cluster.CertificateAuthority = localAuthPath

	// client crt
	localClientCertPath := base.GenLocalCertClientPath(d.remoteAddr)
	err = scp.TransferFile(scp.TransferConfig{
		LocalPath:  localClientCertPath,
		RemoteAddr: d.remoteAddr,
		RemotePath: d.ck.User.ClientCertificate,
	})

	if err != nil {
		return err
	}
	d.ck.User.ClientCertificate = localClientCertPath

	// client key
	localClientKeyPath := base.GenLocalCertClientKeyPath(d.remoteAddr)
	err = scp.TransferFile(scp.TransferConfig{
		LocalPath:  localClientKeyPath,
		RemoteAddr: d.remoteAddr,
		RemotePath: d.ck.User.ClientKey,
	})

	if err != nil {
		return err
	}
	d.ck.User.ClientKey = localClientKeyPath

	return nil
}

// LocalCachePath returns cache path for remote kubectl config by convention.
// that's, `~/.config/cube/cache/$HOST`.
func LocalCachePath(remoteAddr string) string {
	filename := base.ExtractHost(remoteAddr)
	return filepath.Join(base.DefaultCacheDir, filename+".yaml")
}
