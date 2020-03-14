package base

import (
	"fmt"
	"net/url"
	"strings"
)

// SshHost returns host used for ssh/scp connection
func SshHost(user, ip string) string {
	if user == "" {
		return ip
	}

	return user + "@" + ip
}

// ExtractHost extracts host info from remoteAddr which is in the format `user@host`
func ExtractHost(remoteAddr string) string {
	tokens := strings.Split(remoteAddr, "@")

	return tokens[len(tokens)-1]
}

// GetHost returns Host part of address, Host or Host:port if port given.
func GetHost(srvAddr string) string {
	if !strings.HasPrefix(srvAddr, "http") {
		srvAddr = "http://" + srvAddr
	}

	u, err := url.Parse(srvAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to get host:port from address - %v", srvAddr))
	}

	return u.Host
}

// TODO: add tests
// GetHostname returns host, stripping any valid port number if present.
func GetHostname(addr string) string {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	u, err := url.Parse(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to get hostname from address - %v", addr))
	}

	return u.Hostname()
}
