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

func TestGetHost(t *testing.T) {
	tests := []struct {
		name string

		// input
		addr string

		// output
		host string
	}{
		{"hostname-only", "kubernetes", "kubernetes"},
		{"hostname-w/o-port", "https://kubernetes", "kubernetes"},
		{"hostname-w-port", "https://kubernetes:6443", "kubernetes:6443"},
		{"hostname-w/o-schema", "kubernetes:8080", "kubernetes:8080"},
		{"ip-w-port", "https://172.17.1.1:6443", "172.17.1.1:6443"},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			h := getHost(test.addr)
			assert.Equal(test.host, h)
		})
	}
}

func TestGetPort(t *testing.T) {
	tests := []struct {
		name string

		// input
		addr string

		// output
		port int
	}{
		{"hostname-only", "kubernetes", 80},
		{"hostname-w/o-port", "https://kubernetes", 80},
		{"hostname-w-port", "https://kubernetes:6443", 6443},
		{"hostname-w/o-schema", "kubernetes:8080", 8080},
		{"ip-w-port", "https://172.17.1.1:6443", 6443},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			p := getPort(test.addr)
			assert.Equal(test.port, p)
		})
	}
}
