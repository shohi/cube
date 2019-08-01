package kube

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
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
func getNextLocalPort(kc *clientcmdapi.Config) (int, error) {
	// TODO

	return 0, nil
}
