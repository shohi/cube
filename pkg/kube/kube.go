package kube

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/shohi/cube/pkg/config"
	"github.com/shohi/cube/pkg/scp"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// DefaultKubeConfigPath is default kubeconfig path on remote host.
	DefaultKubeConfigPath = "~/.kube/config"

	// DefaultHost is the host used to represent remote master locally
	DefaultHost = "kubernetes"
)

// Fuse merges kubeconfig of remote cluster into local
func Fuse(conf config.Config) error {
	p := getLocalPath(conf.RemoteAddr)

	// TODO: check whether the file is empty
	err := scp.TransferFile(scp.TransferConfig{
		LocalPath:  p,
		RemoteAddr: conf.RemoteAddr,
		RemotePath: DefaultKubeConfigPath,
	})

	if err != nil {
		return err
	}

	// TODO: handle duplicated
	km := newKubeManager(kubeOptions{
		mainPath:   getLocalKubePath(),
		inPath:     p,
		isPurge:    conf.Purge,
		localPort:  conf.LocalPort,
		nameSuffix: conf.NameSuffix,
	})

	if err := km.Do(); err != nil {
		return err
	}

	// TODO: implement dry-run feature
	if conf.DryRun {
		content, err := km.Write()
		if err != nil {
			return err
		}

		log.Printf("merged config:\n%v\n", string(content))
	} else {
		km.WriteToFile()
	}

	// NOTE: Print SSH forwarding setting
	sshCmd := getPortForwardingCmd(km.opts.localPort, km.inAPIAddr, conf.SSHVia)
	log.Printf("ssh command:\n%s\n", sshCmd)

	return nil
}

type kubeOptions struct {
	// TODO: renaming
	mainPath string
	inPath   string

	isPurge    bool
	localPort  int
	nameSuffix string
}

type KubeManager struct {
	opts kubeOptions

	mainKC clientcmdapi.Config
	inKC   clientcmdapi.Config

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

	if k.opts.isPurge {
		return k.purge()
	}

	return k.merge()
}

// merge merges kubeconfig from given path
func (k *KubeManager) merge() error {
	// TODO: remove duplicated
	if len(k.inKC.Clusters) != 1 {
		return errors.New("remote kubeconfig must have only one cluster")
	}

	// NOTE: Only care about `clusters/contexts/users` sections
	for ck, v := range k.inKC.Clusters {
		cluster, found := findCluster(&k.mainKC, v.CertificateAuthorityData)
		if found {
			return fmt.Errorf("kubeconfig already merged under cluster - [%v]", cluster)
		}

		// NOTE: It's ok to set `inAPIAddr` in the loop as
		// only one element is in the map.
		k.inAPIAddr = getHost(v.Server)

		// add Cluster
		ctxName, kCtx := getContext(&k.inKC, ck)
		kUser := getUser(&k.inKC, kCtx.AuthInfo)

		// update Cluster/User/Context
		v.Server = fmt.Sprintf("https://%s:%d", DefaultHost, k.opts.localPort)

		kCtx.AuthInfo += "-" + k.opts.nameSuffix
		kCtx.Cluster += "-" + k.opts.nameSuffix
		ctxName += "-" + k.opts.nameSuffix

		k.mainKC.Clusters[kCtx.Cluster] = v
		k.mainKC.AuthInfos[kCtx.AuthInfo] = kUser
		k.mainKC.Contexts[ctxName] = kCtx
	}

	// TODO: compare based on `certificate-authority-data` field of cluster
	return nil
}

// getHost returns Host part of address, Host or Host:port if port given.
func getHost(srvAddr string) string {
	u, err := url.Parse(srvAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to get host:port from address - %v", srvAddr))
	}

	return u.Host
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

// purge deletes kubeconfig which matches the content with given file path
func (k *KubeManager) purge() error {
	// TODO
	return nil
}

func (k *KubeManager) WriteToFile() error {
	return clientcmd.WriteToFile(k.mainKC, k.opts.mainPath)
}

func (k *KubeManager) Write() ([]byte, error) {
	return clientcmd.Write(k.mainKC)
}
