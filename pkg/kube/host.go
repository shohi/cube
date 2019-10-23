package kube

import (
	"fmt"
	"net"
	"strings"

	"github.com/shohi/cube/pkg/base"
)

const (
	SepAt     = "@"
	SepHyphen = "-"
	SepColon  = ":"

	DefaultRemotePort = 6443
)

// getRemoteHostFromCtx extract host info from kube context.
// e.g full context name - `kubernetes-admin@172.31.7.182:6443-test`,,,,
// result is `172.31.7.182:6443`. If not port provided, `6443` will be used.
func getRemoteHostFromCtx(kctx string) string {
	tokens := strings.Split(kctx, SepAt)
	if len(tokens) < 2 {
		return ""
	}
	remain := strings.Join(tokens[1:], SepAt)
	tokens = strings.Split(remain, SepHyphen)
	if len(tokens) < 2 {
		return ""
	}

	host := tokens[0]
	ipStr := base.GetHostname(host)

	// ipStr is not a IP addr
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	port, notFound := base.GetPort(host)
	if notFound {
		port = DefaultRemotePort
	}

	return fmt.Sprintf("%v:%v", ipStr, port)
}

// getShortContext return context name without user info part.
// e.g full context name - `kubernetes-admin@172.31.7.182:6443-test`
// then short context name is `172.31.7.182-test`.
// NOTE: only remove the part before first `@`.
func getShortContext(kctx string) string {
	tokens := strings.Split(kctx, SepAt)
	if len(tokens) < 2 {
		return kctx
	}

	remain := strings.Join(tokens[1:], SepAt)
	tokens = strings.Split(remain, SepHyphen)

	if len(tokens) < 2 {
		return remain
	}

	hostname := base.GetHostname(tokens[0])
	namesuffix := strings.Join(tokens[1:], SepHyphen)

	return fmt.Sprintf("%v-%v", hostname, namesuffix)
}
