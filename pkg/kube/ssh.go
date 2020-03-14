package kube

import (
	"fmt"
	"os"
)

func GetPortForwardingCmd(localPort int, remoteAPIAddr string, via string) string {
	// example: ssh -fN -L 7002:172.31.6.103:6443 root@xx.xx.xx
	forwardingFmt := "ssh -fN -L %v:%v %v"
	if via == "" {
		if sv := os.Getenv("SSH_VIA"); len(sv) > 0 {
			via = sv
		} else {
			via = "${SSH_VIA}"
		}
	}

	return fmt.Sprintf(forwardingFmt, localPort, remoteAPIAddr, via)
}
