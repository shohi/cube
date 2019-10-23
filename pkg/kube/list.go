package kube

import (
	"fmt"
	"log"
	"strings"

	"github.com/shohi/cube/pkg/base"
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

	return ret, nil
}

func genClusterInfo(kctx string, port int) ClusterInfo {

	info := ClusterInfo{
		Name: getShortContext(kctx),
	}

	addr := getRemoteAddrFromCtx(kctx)
	if addr == "" {
		return info
	}

	info.SSHForward = getPortForwardingCmd(port, addr, "")
	return info
}
