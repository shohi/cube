package action

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

var enableSCPForTest bool

func init() {
	flag.BoolVar(&enableSCPForTest, "enable-scp", false, "whether enable scp in tests")
	// NOTE: DON'T call `Parse` here, which will overwrite default go test flags.
	// flag.Parse()
}

func TestAction_Add(t *testing.T) {
	if !enableSCPForTest {
		t.Skip()
	}

	assert := assert.New(t)

	conf := AddConfig{
		RemoteUser: "core",
		RemoteIP:   "172.31.10.34",
		LocalPort:  7003,
		NameSuffix: "qa",
		DryRun:     true,
	}

	err := Add(conf)
	assert.Nil(err)
}
