package kube

import (
	"testing"

	"github.com/shohi/cube/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestKube_Del(t *testing.T) {
	assert := assert.New(t)

	conf := config.Config{
		RemoteUser: "core",
		RemoteIP:   "172.31.6.103",
		Purge:      true,
		DryRun:     true,
	}

	err := Del(conf)
	assert.Nil(err)
}
