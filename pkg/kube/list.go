package kube

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/shohi/cube/pkg/base"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	errClusterNotFound          = errors.New("cube: cluster not found")
	errLocalPortNotFound        = errors.New("cube: local port not found")
	errLocalServerNotKubernetes = errors.New("cube: local server not 'kubernetes'")
	errNoAddrInContextName      = errors.New("cube: no addr in context name")
)

const (
	defaultRemoteAPIPort = 6443
)

type ClusterInfo struct {
	Name       string `json:"name"`
	SSHForward string `json:"sshForward"`
}

func (f ClusterInfo) String() string {
	return fmt.Sprintf("%s => %s", f.Name, f.SSHForward)
}

type ClusterInfos []ClusterInfo

func (c ClusterInfos) String() string {
	var sb strings.Builder
	for _, v := range c {
		sb.WriteString(v.String())
		sb.WriteByte('\n')
	}

	return sb.String()
}

func (c ClusterInfos) Len() int {
	return len(c)
}

func (c ClusterInfos) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

func (c ClusterInfos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func ListAllClusters() (ClusterInfos, error) {
	kc, err := Load(base.GetLocalKubePath())
	if err != nil {
		return nil, err
	}

	var ret ClusterInfos

	for k := range kc.Contexts {
		info, err := ParseContext(kc, k)
		if err == nil {
			ret = append(ret, *info)
			continue
		}

		if err != nil && !errors.Is(err, errLocalServerNotKubernetes) {
			log.Println(err)
		}
	}

	sort.Sort(ret)

	return ret, nil
}

func ParseContext(kc *clientcmdapi.Config, ctxName string) (*ClusterInfo, error) {
	ctx := kc.Contexts[ctxName]
	cluster, ok := kc.Clusters[ctx.Cluster]
	if !ok {
		return nil, errors.Wrapf(errClusterNotFound, "ctx: %v", ctxName)
	}

	h, p, err := GetOccupiedLocalPort(cluster.Server)
	if err != nil {
		return nil, errors.Wrapf(errLocalPortNotFound, "error %v", err)
	}

	if h != DefaultHost {
		return nil, errors.Wrapf(errLocalServerNotKubernetes, "ctx: %v", ctxName)
	}

	info := genClusterInfo(ctxName, p)
	if info.SSHForward == "" {
		return nil, errors.Wrapf(errNoAddrInContextName, "ctx: %v", ctxName)
	}

	return &info, nil
}

func genClusterInfo(kctx string, port int) ClusterInfo {
	info := ClusterInfo{
		Name: getShortContext(kctx),
	}

	h := getRemoteHostFromCtx(kctx)
	if h == "" {
		return info
	}

	// TODO: dynamicially get real remote port by parsing related kube config file.
	info.SSHForward = GetPortForwardingCmd(port, h, "")
	return info
}
