package huma

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("error")
}

func requireMarshalError(
	t *testing.T,
	f interface {
		Marshal(writer io.Writer, input any) error
	},
) {
	t.Helper()

	err := f.Marshal(failingWriter{}, "x")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMarshal)
}

func requireUnmarshalError(
	t *testing.T,
	f interface {
		Unmarshal(input []byte, result any) error
	},
) {
	t.Helper()

	var out any

	err := f.Unmarshal([]byte("invalid"), &out)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnmarshal)
}

func TestFormatTextPlain(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", `"hello"`},
		{"int", 42, `42`},
		{"struct", struct{ A string }{"value"}, `{"A":"value"}`},
	}

	frmt := NewFormatTextPlain()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := frmt.Marshal(&buf, tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}

	t.Run("marshal error", func(t *testing.T) {
		requireMarshalError(t, frmt)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		requireUnmarshalError(t, frmt)
	})
}

func TestFormatTextHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string",
			input:    "hello",
			expected: "<!doctype html>\n<html>\n  <body>\n    <pre>&#34;hello&#34;</pre>\n  </body>\n</html>\n",
		},
		{
			name:     "angle brackets",
			input:    "<tag>",
			expected: "<!doctype html>\n<html>\n  <body>\n    <pre>&#34;&lt;tag&gt;&#34;</pre>\n  </body>\n</html>\n",
		},
		{
			name:     "struct",
			input:    struct{ Name string }{"bob"},
			expected: "<!doctype html>\n<html>\n  <body>\n    <pre>{&#34;Name&#34;:&#34;bob&#34;}</pre>\n  </body>\n</html>\n",
		},
		{
			name:  "special chars",
			input: `"quote" & 'apostrophe'`,
			expected: `<!doctype html>
<html>
  <body>
    <pre>&#34;\&#34;quote\&#34; &amp; &#39;apostrophe&#39;&#34;</pre>
  </body>
</html>
`,
		},
	}

	frmt := NewFormatTextHTML()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := frmt.Marshal(&buf, tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}

	t.Run("marshal error", func(t *testing.T) {
		requireMarshalError(t, frmt)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		requireUnmarshalError(t, frmt)
	})
}

func TestFormatTextHTML_WithFormatTextHTMLTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string",
			input:    "hello",
			expected: "<!doctype html>\n<html>\n  <head>test</head>\n  <body>\n    <pre>&#34;hello&#34;</pre>\n  </body>\n</html>\n",
		},
		{
			name:     "angle brackets",
			input:    "<tag>",
			expected: "<!doctype html>\n<html>\n  <head>test</head>\n  <body>\n    <pre>&#34;&lt;tag&gt;&#34;</pre>\n  </body>\n</html>\n",
		},
		{
			name:     "struct",
			input:    struct{ Name string }{"bob"},
			expected: "<!doctype html>\n<html>\n  <head>test</head>\n  <body>\n    <pre>{&#34;Name&#34;:&#34;bob&#34;}</pre>\n  </body>\n</html>\n",
		},
	}

	fmt := NewFormatTextHTML(
		WithFormatTextHTMLTemplate(`<!doctype html>
<html>
  <head>test</head>
  <body>
    <pre>{{ . }}</pre>
  </body>
</html>
`),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := fmt.Marshal(&buf, tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestFormatApplicationJSON(t *testing.T) {
	type sample struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name  string
		input sample
	}{
		{"simple struct", sample{Name: "alice", Age: 30}},
	}

	frmt := NewFormatApplicationJSON()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := frmt.Marshal(&buf, tt.input)
			require.NoError(t, err)

			var out sample

			err = frmt.Unmarshal(buf.Bytes(), &out)
			require.NoError(t, err)

			require.Equal(t, tt.input, out)
		})
	}

	t.Run("marshal error", func(t *testing.T) {
		requireMarshalError(t, frmt)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		var out any

		err := frmt.Unmarshal([]byte("{invalid"), &out)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnmarshal)
	})
}

func TestFormat_IsMarshalIsUnmarshal(t *testing.T) {
	errBoom := errors.New("error")

	formatters := map[string]interface {
		IsMarshal(err error) bool
		IsUnmarshal(err error) bool
	}{
		"FormatTextPlain":       NewFormatTextPlain(),
		"FormatTextHTML":        NewFormatTextHTML(),
		"FormatApplicationJSON": NewFormatApplicationJSON(),
	}

	for name, frmt := range formatters {
		t.Run(name, func(t *testing.T) {
			require.True(t, frmt.IsMarshal(ErrMarshal), "bare marshal sentinel")
			require.True(
				t,
				frmt.IsMarshal(fmt.Errorf("%w: %w", ErrMarshal, errBoom)),
				"wrapped marshal sentinel",
			)
			require.False(t, frmt.IsMarshal(errBoom), "unrelated error")
			require.False(t, frmt.IsMarshal(ErrUnmarshal), "the other sentinel")
			require.False(t, frmt.IsMarshal(nil), "nil error")

			require.True(t, frmt.IsUnmarshal(ErrUnmarshal), "bare unmarshal sentinel")
			require.True(
				t,
				frmt.IsUnmarshal(fmt.Errorf("%w: %w", ErrUnmarshal, errBoom)),
				"wrapped unmarshal sentinel",
			)
			require.False(t, frmt.IsUnmarshal(errBoom), "unrelated error")
			require.False(t, frmt.IsUnmarshal(ErrMarshal), "the other sentinel")
			require.False(t, frmt.IsUnmarshal(nil), "nil error")
		})
	}
}
