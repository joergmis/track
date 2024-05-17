package track_test

import (
	"testing"

	"github.com/joergmis/track"
	"github.com/stretchr/testify/assert"
)

func TestBackend_Valid(t *testing.T) {
	tests := []struct {
		name    string
		backend string
		valid   bool
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(
				t,
				track.BackendType(test.backend),
				test.valid,
			)
		})
	}
}
