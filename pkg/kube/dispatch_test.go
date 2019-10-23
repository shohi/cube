package kube

import (
	"flag"
	"testing"

	"github.com/shohi/cube/pkg/config"
	"github.com/stretchr/testify/assert"
)

var enableSCPForTest bool

func init() {
	flag.BoolVar(&enableSCPForTest, "enable-scp", false, "whether enable scp in tests")
	// NOTE: DON'T call `Parse` here, which will overwrite default go test flags.
	// flag.Parse()
}

func TestDispatch_Merge(t *testing.T) {
	if !enableSCPForTest {
		t.Skip()
	}

	assert := assert.New(t)

	conf := config.Config{
		RemoteUser: "core",
		RemoteIP:   "172.31.10.34",
		LocalPort:  7003,
		NameSuffix: "qa",
		DryRun:     true,
	}

	err := Dispatch(conf)
	assert.Nil(err)
}

func TestDispatch_Purge(t *testing.T) {
	if !enableSCPForTest {
		t.Skip()
	}

	assert := assert.New(t)

	conf := config.Config{
		RemoteUser: "core",
		RemoteIP:   "172.31.6.103",
		Purge:      true,
		DryRun:     true,
	}

	err := Dispatch(conf)

	assert.Nil(err)
}
