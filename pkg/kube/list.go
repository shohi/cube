package kube

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/shohi/cube/pkg/base"
)

const (
	defaultRemoteAPIPort = 6443
)

type ClusterInfo struct {
	Name       string
	SSHForward string
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
	kc, err := load(base.GetLocalKubePath())
	if err != nil {
		return nil, err
	}

	var ret ClusterInfos

	for k, v := range kc.Clusters {
		h, p, err := getOccupiedLocalPort(v.Server)
		if err != nil {
			// TODO: use debug level
			log.Printf("failed to get server port for cluster - [%v], err: %v", k, err)
			continue
		}

		if h != DefaultHost {
			continue
		}

		ctxName, _ := getContext(kc, k)
		info := genClusterInfo(ctxName, p)
		if info.SSHForward == "" {
			// TODO: use debug level
			log.Printf("no remote addr found in context name for cluster - [%v]", k)
			continue
		}
		ret = append(ret, info)
	}

	sort.Sort(ret)

	return ret, nil
}

func genClusterInfo(kctx string, port int) ClusterInfo {

	info := ClusterInfo{
		Name: getShortContext(kctx),
	}

	ip := getRemoteIPFromCtx(kctx)
	if ip == "" {
		return info
	}

	// TODO: dynamicially get real remote port by parsing related kube config file.
	addr := fmt.Sprintf("%v:%v", ip, defaultRemoteAPIPort)
	info.SSHForward = getPortForwardingCmd(port, addr, "")
	return info
}
