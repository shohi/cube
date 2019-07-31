package kube

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPortForwardingCmd(t *testing.T) {

	localPort := 7001
	remoteAPIAddr := "172.17.31.1:6443"

	tests := []struct {
		name string

		// input
		viaEnv string

		// output
		expResult string
	}{
		{"env-exist",
			"core@192.168.1.1", "ssh -fN -L 7001:172.17.31.1:6443 core@192.168.1.1"},
		{"env-nonexist",
			"", "ssh -fN -L 7001:172.17.31.1:6443 ${SSH_VIA}"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			os.Setenv("SSH_VIA", test.viaEnv)
			ret := getPortForwardingCmd(localPort, remoteAPIAddr, test.viaEnv)

			assert.Equal(test.expResult, ret)
		})

	}

}
