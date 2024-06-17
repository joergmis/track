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
	}{
		{
			name:    "typo",
			backend: "l0cal",
			valid:   false,
		},
		{
			name:    "non-existend",
			backend: "some-backend",
			valid:   false,
		},
		{
			name:    "local",
			backend: string(track.BackendLocal),
			valid:   true,
		},
		{
			name:    "clockodo",
			backend: string(track.BackendClockodo),
			valid:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(
				t,
				track.BackendType(test.backend).Valid(),
				test.valid,
			)
		})
	}
}
