package kube

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	minLocalPort = 7001
	maxLocalPort = 7100

	defaultHost = "kubernetes"
)

// extractHost extracts host info from remoteAddr which is in the format `user@host`
func extractHost(remoteAddr string) string {
	tokens := strings.Split(remoteAddr, "@")

	return tokens[len(tokens)-1]
}

// getHost returns Host part of address, Host or Host:port if port given.
func getHost(srvAddr string) string {
	if !strings.HasPrefix(srvAddr, "http") {
		srvAddr = "http://" + srvAddr
	}

	u, err := url.Parse(srvAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to get host:port from address - %v", srvAddr))
	}

	return u.Host
}

// getPort returns port part of address.
func getPort(srvAddr string) int {
	if !strings.HasPrefix(srvAddr, "http") {
		srvAddr = "http://" + srvAddr
	}

	u, err := url.Parse(srvAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to get port from address - %v", srvAddr))
	}

	p := u.Port()
	if p == "" {
		p = "80"
	}

	pVal, err := strconv.Atoi(p)
	if err != nil {
		panic(err)
	}

	return pVal
}

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
		if isAvailable(k) {
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
		u, err := url.Parse(v.Server)
		if err != nil {
			log.Printf("failed to parse server address - [%v]\n", v.Server)
			continue
		}

		portStr := u.Port()
		if portStr == "" {
			continue
		}

		if p, err := strconv.ParseInt(portStr, 10, 32); err == nil {
			ret = append(ret, int(p))
		}
	}

	return ret
}

// isAvailable tests whether given port is available on localhost.
func isAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}

	defer func() {
		if ln != nil {
			_ = ln.Close()
		}
	}()

	return true
}
