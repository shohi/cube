package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHost(t *testing.T) {
	tests := []struct {
		name string

		// input
		addr string

		// output
		host string
	}{
		{"ip-w/o-at", "172.17.1.1", "172.17.1.1"},
		{"host-w/o-at", "kubernetes", "kubernetes"},
		{"ip-w-at", "user@172.17.1.1", "172.17.1.1"},
		{"host-w-at", "core@kubernetes", "kubernetes"},
		{"empty", "", ""},
		{"multiple-at", "user@pass@host", "host"},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			h := extractHost(test.addr)
			assert.Equal(test.host, h)
		})
	}

}
