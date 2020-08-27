package kube

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/shohi/cube/pkg/base"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// DefaultHost is the host used to represent remote master locally
	DefaultHost = "kubernetes"
)

func getRemoteAddr(user, ip string) string {
	if user == "" {
		return ip
	}

	return user + "@" + ip
}

func WriteToFile(kc *clientcmdapi.Config, configPath string) error {
	return clientcmd.WriteToFile(*kc, configPath)
}

func Write(kc *clientcmdapi.Config) ([]byte, error) {
	return clientcmd.Write(*kc)
}

func findCluster(kc *clientcmdapi.Config, inCluster *clientcmdapi.Cluster) (string, bool) {
	var cluster string
	var found = false

	for k, v := range kc.Clusters {
		if isEqual(v, inCluster) {
			found = true
			cluster = k
			break
		}
	}

	return cluster, found
}

// isEqual checks whether the two cluster is equal based on CertificateAuthority info.
func isEqual(k1, k2 *clientcmdapi.Cluster) bool {
	if len(k1.CertificateAuthorityData) > 0 &&
		string(k1.CertificateAuthorityData) == string(k2.CertificateAuthorityData) {
		return true
	}

	if len(k1.CertificateAuthority) > 0 &&
		k1.CertificateAuthority == k2.CertificateAuthority {
		return true
	}

	return false
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

func FindContextsByName(kc *clientcmdapi.Config, name string, filter func(string) bool) map[string]*clientcmdapi.Context {
	var result = make(map[string]*clientcmdapi.Context)
	for k, v := range kc.Contexts {
		if !strings.Contains(k, name) {
			continue
		}

		if filter != nil && !filter(k) {
			continue
		}

		result[k] = v
	}

	return result
}

// Load reads kubeconfig from file
func Load(configPath string) (*clientcmdapi.Config, error) {
	exist, isDir := FileExists(configPath)

	if !exist {
		conf := clientcmdapi.NewConfig()
		return conf, nil
	} else if isDir {
		return nil, fmt.Errorf("config path [%v] is a dir, not file", configPath)
	}

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

// remote API address is composed of ip from remoteAddr and port for apiSrv.
// e.g. 172.10.0.1:6443
func genRemoteAPIAddr(remoteAddr, apiSrv string) string {
	h := base.GetHostname(remoteAddr)
	p, _ := base.GetPort(apiSrv)
	return fmt.Sprintf("%v:%v", h, p)
}
