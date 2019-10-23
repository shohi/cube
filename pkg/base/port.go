package base

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// GetPort returns port part of address.
func GetPort(srvAddr string) (port int, notFound bool) {
	if !strings.HasPrefix(srvAddr, "http") {
		srvAddr = "http://" + srvAddr
	}

	u, err := url.Parse(srvAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to get port from address - %v", srvAddr))
	}

	p := u.Port()
	if p == "" {
		if u.Scheme == "https" {
			return 443, true
		}

		return 80, true
	}

	pVal, err := strconv.Atoi(p)
	if err != nil {
		panic(err)
	}

	return pVal, false
}

// IsAvailable tests whether given port is available on localhost.
func IsAvailable(port int) bool {
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
