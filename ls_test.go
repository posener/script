package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		paths     []string
		wantError bool
		want      string
	}{
		{
			name:  "single file",
			paths: []string{"testdata/a.txt"},
			want:  "testdata/a.txt\n",
		},
		{
			name:  "directory",
			paths: []string{"testdata"},
			want:  "testdata/a.txt\ntestdata/b.txt\n",
		},
		{
			name:  "multiple paths",
			paths: []string{"testdata", "testdata/a.txt"},
			want:  "testdata/a.txt\ntestdata/b.txt\ntestdata/a.txt\n",
		},
		{
			name:      "error",
			paths:     []string{"no-such-file"},
			wantError: true,
			want:      "",
		},
		{
			name:      "error with successful path",
			paths:     []string{"no-such-file", "testdata/a.txt"},
			wantError: true,
			want:      "testdata/a.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ls(tt.paths...).ToString()
			require.Equal(t, tt.wantError, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
