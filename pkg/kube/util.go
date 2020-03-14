package kube

import (
	"fmt"

	"github.com/shohi/cube/pkg/base"
)

// remote API address is composed of ip from remoteAddr and port for apiSrv.
// e.g. 172.10.0.1:6443
func genInAPIAddr(remoteAddr, apiSrv string) string {
	h := base.GetHostname(remoteAddr)
	p, _ := base.GetPort(apiSrv)
	return fmt.Sprintf("%v:%v", h, p)
}
