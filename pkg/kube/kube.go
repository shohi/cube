package kube

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// DefaultKubeConfigPath is default kubeconfig path on remote host.
	DefaultKubeConfigPath = "~/.kube/config"

	// DefaultHost is the host used to represent remote master locally
	DefaultHost = "kubernetes"
)

var (
	ErrConfigAlreadyMerged = errors.New("kubeconfig has already been merged")
	ErrConfigInvalid       = errors.New("remote kubeconfig must have only one cluster")

	ErrInvalidLocalPort = errors.New("invalid local port for merging")
	ErrEmptyNameSuffix  = errors.New("name-suffix must not be empty for merging")
)

func getRemoteAddr(user, ip string) string {
	if user == "" {
		return ip
	}

	return user + "@" + ip
}

type kubeOptions struct {
	// TODO: renaming
	mainPath string
	inPath   string

	action     ActionType
	localPort  int
	nameSuffix string
}

type KubeManager struct {
	opts kubeOptions

	mainKC clientcmdapi.Config
	inKC   clientcmdapi.Config

	inClusterName string
	inCluster     *clientcmdapi.Cluster

	inCtxName string
	inCtx     *clientcmdapi.Context

	inUser *clientcmdapi.AuthInfo

	updatedClusterName string

	inAPIAddr string
}

func newKubeManager(opts kubeOptions) *KubeManager {
	return &KubeManager{opts: opts}
}

func (k *KubeManager) init() error {
	kc, err := load(k.opts.mainPath)
	if err != nil {
		return err
	}
	k.mainKC = *kc

	kc, err = load(k.opts.inPath)
	if err != nil {
		return err
	}
	k.inKC = *kc

	return nil
}

func (k *KubeManager) Do() error {
	if err := k.init(); err != nil {
		return err
	}

	if err := k.extractInKC(); err != nil {
		return err
	}

	switch k.opts.action {
	case ActionPurge:
		return k.purge()
	case ActionPrint:
		return k.inferLocalPort()
	default:
		return k.merge()
	}
}

// extractInKC extracts Cluster/User/Context info from `inKC`.
func (k *KubeManager) extractInKC() error {
	// TODO: remove duplicated
	if len(k.inKC.Clusters) != 1 {
		return ErrConfigInvalid
	}

	// NOTE: Only care about `clusters/contexts/users` sections
	for ck, v := range k.inKC.Clusters {
		// Take a snapshot of incoming cluster info
		// NOTE: It's ok to set `inAPIAddr` in the loop as
		// only one element is in the map.
		k.inClusterName, k.inCluster = ck, v
		k.inCtxName, k.inCtx = getContext(&k.inKC, ck)
		k.inUser = getUser(&k.inKC, k.inCtx.AuthInfo)

		k.inAPIAddr = getHost(v.Server)
	}

	// TODO: compare based on `certificate-authority-data` field of cluster
	return nil
}

func (k *KubeManager) merge() error {
	// check
	if k.opts.localPort == 0 {
		k.opts.localPort = getNextLocalPort(&k.mainKC)
	}

	if k.opts.localPort <= 0 || k.opts.localPort > 65535 {
		return errors.Wrapf(ErrInvalidLocalPort, "local port: [%v]", k.opts.localPort)
	}

	if k.opts.nameSuffix == "" {
		return ErrEmptyNameSuffix
	}

	cluster, found := findCluster(&k.mainKC, k.inCluster.CertificateAuthorityData)
	if found {
		return errors.Wrapf(ErrConfigAlreadyMerged, "cluster: [%v]", cluster)
	}

	// add Cluster
	k.inCluster.Server = fmt.Sprintf("https://%s:%d", DefaultHost, k.opts.localPort)

	k.inCtx.AuthInfo += "-" + k.opts.nameSuffix
	k.inCtx.Cluster += "-" + k.opts.nameSuffix
	k.inCtxName += "-" + k.opts.nameSuffix

	k.updatedClusterName = k.inCtx.Cluster

	if _, ok := k.mainKC.Clusters[k.inCtx.Cluster]; ok {
		return fmt.Errorf("cluster - [%v] - already exists, plz choose another suffix", k.inCtx.Cluster)
	}

	if _, ok := k.mainKC.AuthInfos[k.inCtx.AuthInfo]; ok {
		return fmt.Errorf("user - [%v] - already exists, plz choose anthor suffix", k.inCtx.AuthInfo)
	}

	if _, ok := k.mainKC.Contexts[k.inCtxName]; ok {
		return fmt.Errorf("context - [%v] - already exists, plz choose another suffix", k.inCtxName)
	}

	k.mainKC.Clusters[k.inCtx.Cluster] = k.inCluster
	k.mainKC.AuthInfos[k.inCtx.AuthInfo] = k.inUser
	k.mainKC.Contexts[k.inCtxName] = k.inCtx

	return nil
}

// purge deletes kubeconfig which matches the content with given file path
func (k *KubeManager) purge() error {
	cluster, found := findCluster(&k.mainKC, k.inCluster.CertificateAuthorityData)
	if !found {
		return fmt.Errorf("cannot find matched kubeconfig for purging")
	}

	// update local port from kubeconfig
	c := k.mainKC.Clusters[cluster]
	k.opts.localPort = getPort(c.Server)

	k.updatedClusterName = cluster

	// TODO: move to other place
	fmt.Fprintf(os.Stdout, "# cluster - [%v] - will be purged\n", cluster)

	ctxName, ctx := getContext(&k.mainKC, cluster)

	delete(k.mainKC.Clusters, ctx.Cluster)
	delete(k.mainKC.AuthInfos, ctx.AuthInfo)
	delete(k.mainKC.Contexts, ctxName)

	return nil
}

// inferLocalPort infers local port from default kubeconfig if local-port not provided.
func (k *KubeManager) inferLocalPort() error {
	if k.opts.localPort > 0 && k.opts.localPort < 65535 {
		return nil
	}

	cluster, found := findCluster(&k.mainKC, k.inCluster.CertificateAuthorityData)
	if !found {
		return errors.Wrapf(ErrInvalidLocalPort, "local port: [%v]", k.opts.localPort)
	}

	c := k.mainKC.Clusters[cluster]
	k.opts.localPort = getPort(c.Server)

	return nil
}

func (k *KubeManager) WriteToFile() error {
	return clientcmd.WriteToFile(k.mainKC, k.opts.mainPath)
}

func (k *KubeManager) Write() ([]byte, error) {
	return clientcmd.Write(k.mainKC)
}

func findCluster(kc *clientcmdapi.Config, certAuthData []byte) (string, bool) {
	var cluster string
	var found = false

	for k, v := range kc.Clusters {
		if string(v.CertificateAuthorityData) == string(certAuthData) {
			found = true
			cluster = k
			break
		}
	}

	return cluster, found
}

func getContext(kc *clientcmdapi.Config, cluster string) (string, *clientcmdapi.Context) {
	var kCtx *clientcmdapi.Context
	var kName string

	var found = false

	for k, v := range kc.Contexts {
		if v.Cluster == cluster {
			found = true
			kName = k
			kCtx = v
			break
		}
	}

	if !found {
		panic(fmt.Sprintf("failed to extract context info with given cluster - %v", cluster))
	}

	return kName, kCtx
}

func getUser(kc *clientcmdapi.Config, userName string) *clientcmdapi.AuthInfo {
	user, ok := kc.AuthInfos[userName]
	if !ok {
		panic(fmt.Sprintf("failed to extract auth info with given user - %v", userName))
	}

	return user
}
