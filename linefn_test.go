package script

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testModifier(in []byte) ([]byte, error) {
	// In case of EOF:
	if in == nil {
		return nil, nil
	}

	in = append([]byte{'@'}, in...)
	return append(in, '@', '\n'), nil
}

func testEOFModifier(in []byte) ([]byte, error) {
	return nil, io.EOF
}

func testErrorModifier(in []byte) ([]byte, error) {
	return nil, fmt.Errorf("error")
}

func TestLineFn(t *testing.T) {
	t.Parallel()
	// Create line that is long enough such that it won't be read in a single bufio read-line
	// of 4096 bytes.
	longLine := strings.Repeat("a", 10000)

	tests := []struct {
		name     string
		input    string
		modifier Modifer
		want     string
	}{
		{
			name:     "simple",
			modifier: ModifierFn(testModifier),
			input:    "a\nb\nc",
			want:     "@a@\n@b@\n@c@\n",
		},
		{
			name:     "long line correctness",
			modifier: ModifierFn(testModifier),
			input:    longLine + "\n" + longLine,
			want:     "@" + longLine + "@\n@" + longLine + "@\n",
		},
		{
			name:     "eof handling",
			modifier: ModifierFn(testEOFModifier),
			input:    "a",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Echo(tt.input).LineFn(tt.modifier).ToString()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLineFn_error(t *testing.T) {
	t.Parallel()
	got, err := Echo("a").LineFn(ModifierFn(testErrorModifier)).ToString()
	assert.Error(t, err)
	assert.Equal(t, "", got)
}
