package scp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasContent(t *testing.T) {
	tests := []struct {
		name string

		// input
		filename string

		// output
		expNonEmpty bool
	}{
		{"nonempty", "data.txt", true},
		{"empty", "empty.txt", false},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			ret := hasContent(filepath.Join("testdata", test.filename))
			assert.Equal(test.expNonEmpty, ret)
		})
	}
}
