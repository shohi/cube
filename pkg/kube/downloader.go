package kube

import (
	"errors"
	"path/filepath"

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
	kc         *clientcmdapi.Config

	clusterName string
	cluster     *clientcmdapi.Cluster

	ctxName string
	ctx     *clientcmdapi.Context

	user *clientcmdapi.AuthInfo
}

// NewDownloader create a new remote config downloader.
func NewDownloader(remoteAddr string) *Downloader {
	return &Downloader{
		remoteAddr: remoteAddr,
	}
}

// Download fetches config and cert files.
func (d *Downloader) Download() (*clientcmdapi.Config, error) {
	if err := d.downloadK8sConfig(); err != nil {
		return nil, err
	}

	if err := d.checkCertFiles(); err != nil {
		return nil, err
	}

	return d.kc, nil

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

	d.kc, err = Load(p)
	return err
}

func (d *Downloader) checkCertFiles() error {
	if d.kc == nil {
		return ErrConfigInvalid
	}

	// TODO: remove duplicated
	if len(d.kc.Clusters) != 1 {
		return ErrConfigInvalid
	}

	// NOTE: Only care about `clusters/contexts/users` sections
	for ck, v := range d.kc.Clusters {
		// Take a snapshot of incoming cluster info
		// NOTE: It's ok to set `inAPIAddr` in the loop as
		// only one element is in the map.
		d.clusterName, d.cluster = ck, v
		d.ctxName, d.ctx = getContext(d.kc, ck)
		d.user = getUser(d.kc, d.ctx.AuthInfo)
	}

	// TODO: use logrus
	// log.Printf("=====> cluster: [%v]\n", d.cluster)

	// check if token set. If token is set, no cert file/data is needed.
	if len(d.user.Token) > 0 {
		return nil
	}

	// k8s > 1.7
	if len(d.cluster.CertificateAuthorityData) > 0 {
		if len(d.user.ClientCertificateData) == 0 ||
			len(d.user.ClientKeyData) == 0 {
			return ErrRemoteInvalidUser
		}
		return nil
	}

	// k8s <= 1.7
	if len(d.cluster.CertificateAuthority) > 0 {
		// Download cert files
		return d.downloadCertFiles()
	}

	return ErrRemoteInvalidCert
}

func (d *Downloader) downloadCertFiles() error {
	if len(d.cluster.CertificateAuthority) == 0 {
		return nil
	}

	// download auth cert and also update corresponding info
	localAuthPath := base.GenLocalCertAuthPath(d.remoteAddr)
	err := scp.TransferFile(scp.TransferConfig{
		LocalPath:  localAuthPath,
		RemoteAddr: d.remoteAddr,
		RemotePath: d.cluster.CertificateAuthority,
	})

	if err != nil {
		return err
	}
	d.cluster.CertificateAuthority = localAuthPath

	// client crt
	localClientCertPath := base.GenLocalCertClientPath(d.remoteAddr)
	err = scp.TransferFile(scp.TransferConfig{
		LocalPath:  localClientCertPath,
		RemoteAddr: d.remoteAddr,
		RemotePath: d.user.ClientCertificate,
	})

	if err != nil {
		return err
	}
	d.user.ClientCertificate = localClientCertPath

	// client key
	localClientKeyPath := base.GenLocalCertClientKeyPath(d.remoteAddr)
	err = scp.TransferFile(scp.TransferConfig{
		LocalPath:  localClientKeyPath,
		RemoteAddr: d.remoteAddr,
		RemotePath: d.user.ClientKey,
	})

	if err != nil {
		return err
	}
	d.user.ClientKey = localClientKeyPath

	return nil
}

// LocalCachePath returns cache path for remote kubectl config by convention.
// that's, `~/.config/cube/cache/$HOST`.
func LocalCachePath(remoteAddr string) string {
	filename := base.ExtractHost(remoteAddr)
	return filepath.Join(base.DefaultCacheDir, filename+".yaml")
}
