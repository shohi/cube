package kube

import (
	"log"
	"net/url"
	"strconv"

	"github.com/shohi/cube/pkg/base"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	minLocalPort = 7001
	maxLocalPort = 7100

	defaultHost = "kubernetes"
)

// getNextLocalPort get next available local port.
// It checks the cluster whose server is in format `https://kubernetes:xxx`.
func getNextLocalPort(kc *clientcmdapi.Config) int {
	ocPorts := getAllOccupiedLocalPort(kc)
	var max = minLocalPort - 1

	for _, v := range ocPorts {
		if v < minLocalPort || v > maxLocalPort {
			continue
		}
		if v > max {
			max = v
		}
	}

	for k := max + 1; k < maxLocalPort; k++ {
		if base.IsAvailable(k) {
			return k
		}

	}

	return -1
}

func getAllOccupiedLocalPort(kc *clientcmdapi.Config) []int {
	if kc == nil || len(kc.Clusters) == 0 {
		return nil
	}

	var ret []int

	for _, v := range kc.Clusters {
		if v == nil {
			continue
		}

		_, p, err := GetOccupiedLocalPort(v.Server)
		if err != nil {
			log.Printf("failed to get port from server address - [%v], err: %v\n", v.Server, err)
			continue
		}
		ret = append(ret, p)
	}

	return ret
}

func GetOccupiedLocalPort(srv string) (host string, port int, err error) {
	u, err := url.Parse(srv)
	if err != nil {
		return "", -1, err
	}

	p, err := getPort(u)

	return u.Hostname(), p, err
}

// getPort retrieves port info from url.URL.
func getPort(u *url.URL) (int, error) {
	portStr := u.Port()
	if portStr == "" {
		return 80, nil
	}

	p, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return -1, err
	}

	return int(p), nil
}
