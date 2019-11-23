package script

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadTail(t *testing.T) {
	t.Parallel()

	const text = "a\nbb"
	const text2 = "a\r\nbb"

	tests := []struct {
		src  string
		n    int
		head string
		tail string
	}{
		{src: text, n: -3, head: "", tail: ""},
		{src: text, n: -2, head: "", tail: ""},
		{src: text, n: -1, head: "bb\n", tail: "a\n"},
		{src: text, n: 0, head: "", tail: ""},
		{src: text, n: 1, head: "a\n", tail: "bb\n"},
		{src: text, n: 2, head: "a\nbb\n", tail: "a\nbb\n"},
		{src: text, n: 3, head: "a\nbb\n", tail: "a\nbb\n"},

		{src: text2, n: -3, head: "", tail: ""},
		{src: text2, n: -2, head: "", tail: ""},
		{src: text2, n: -1, head: "bb\n", tail: "a\n"},
		{src: text2, n: 0, head: "", tail: ""},
		{src: text2, n: 1, head: "a\n", tail: "bb\n"},
		{src: text2, n: 2, head: "a\nbb\n", tail: "a\nbb\n"},
		{src: text2, n: 3, head: "a\nbb\n", tail: "a\nbb\n"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("head/%s/%d", tt.src, tt.n), func(t *testing.T) {
			got, err := Echo(tt.src).Head(tt.n).ToString()
			require.NoError(t, err)
			assert.Equal(t, tt.head, got)
		})
		t.Run(fmt.Sprintf("tail/%s/%d", tt.src, tt.n), func(t *testing.T) {
			got, err := Echo(tt.src).Tail(tt.n).ToString()
			require.NoError(t, err)
			assert.Equal(t, tt.tail, got)
		})
	}
}
