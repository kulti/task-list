package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		head string
		tail string
	}{
		{"", "", "/"},
		{"/", "", "/"},
		{"/path", "path", "/"},
		{"/path/", "path", "/"},
		{"/path/path2", "path", "/path2"},
		{"/path/path2/", "path", "/path2"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			head, tail := shiftPath(tc.path)
			assert.Equal(t, tc.head, head)
			assert.Equal(t, tc.tail, tail)
		})
	}
}
