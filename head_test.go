package script

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHead(t *testing.T) {
	t.Parallel()

	const text = "a\nbb"
	const text2 = "a\r\nbb"

	tests := []struct {
		src  string
		n    int
		want string
	}{
		{src: text, n: -3, want: "a\nbb\n"},
		{src: text, n: -2, want: "a\nbb\n"},
		{src: text, n: -1, want: "bb\n"},
		{src: text, n: 0, want: ""},
		{src: text, n: 1, want: "a\n"},
		{src: text, n: 2, want: "a\nbb\n"},
		{src: text, n: 3, want: "a\nbb\n"},

		{src: text2, n: -3, want: "a\nbb\n"},
		{src: text2, n: -2, want: "a\nbb\n"},
		{src: text2, n: -1, want: "bb\n"},
		{src: text2, n: 0, want: ""},
		{src: text2, n: 1, want: "a\n"},
		{src: text2, n: 2, want: "a\nbb\n"},
		{src: text2, n: 3, want: "a\nbb\n"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%d", tt.src, tt.n), func(t *testing.T) {
			got, err := Echo(tt.src).Head(tt.n).ToString()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
