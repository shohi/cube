package base

import (
	"fmt"
	"net/url"
	"strings"
)

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
