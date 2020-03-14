package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAction_Del(t *testing.T) {
	t.Skip("skip")
	assert := assert.New(t)

	conf := DelConfig{
		Name:   "test",
		DryRun: true,
	}

	err := Del(conf)
	assert.Nil(err)
}
