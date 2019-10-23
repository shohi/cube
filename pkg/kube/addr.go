package kube

import (
	"net"
	"strings"
)

const (
	SepAt     = "@"
	SepHyphen = "-"
)

// getRemoteAddrFromCtx extract IP from kube context.
// e.g full context name - `kubernetes-admin@172.31.7.182-test`,,,,
// result is `172.31.7.181`.
func getRemoteAddrFromCtx(kctx string) string {
	sCtx := getShortContext(kctx)
	var ipStr string
	tokens := strings.Split(sCtx, SepHyphen)
	if len(tokens) < 2 {
		ipStr = sCtx
	} else {
		ipStr = tokens[0]
	}

	ip := net.ParseIP(ipStr)
	if ip != nil {
		return ipStr
	}

	return ""
}

// getShortContext return context name without user info part.
// e.g full context name - `kubernetes-admin@172.31.7.182-test`
// then short context name is `172.31.7.182-test`.
// NOTE: only remove the part before first `@`.
func getShortContext(kctx string) string {
	tokens := strings.Split(kctx, SepAt)
	if len(tokens) < 2 {
		return kctx
	}
	return strings.Join(tokens[1:], SepAt)
}
