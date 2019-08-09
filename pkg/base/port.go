package base

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// GetPort returns port part of address.
func GetPort(srvAddr string) int {
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
