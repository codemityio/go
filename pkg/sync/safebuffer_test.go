package sync

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestSafeBuffer_WriteAndString(t *testing.T) {
	sb := NewSafeBuffer()
	n, err := sb.Write([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, "hello", sb.String())
}

func TestSafeBuffer_WriteString(t *testing.T) {
	sb := NewSafeBuffer()
	n, err := sb.WriteString("world")
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, "world", sb.String())
}

func TestSafeBuffer_WriteByte(t *testing.T) {
	sb := NewSafeBuffer()
	err := sb.WriteByte('!')
	require.NoError(t, err)
	require.Equal(t, "!", sb.String())
}

func TestSafeBuffer_WriteRune(t *testing.T) {
	sb := NewSafeBuffer()
	n, err := sb.WriteRune('🎉')
	require.NoError(t, err)
	require.Equal(t, utf8.RuneLen('🎉'), n)
	require.Equal(t, "🎉", sb.String())
}

func TestSafeBuffer_Bytes(t *testing.T) {
	sb := NewSafeBuffer()
	_, _ = sb.WriteString("buffer test")
	bytes := sb.Bytes()
	require.Equal(t, []byte("buffer test"), bytes)
}

func TestSafeBuffer_LenAndCap(t *testing.T) {
	sb := NewSafeBuffer()
	_, _ = sb.WriteString("1234567890")
	require.GreaterOrEqual(t, sb.Cap(), sb.Len())
	require.Equal(t, 10, sb.Len())
}

func TestSafeBuffer_Reset(t *testing.T) {
	sb := NewSafeBuffer()
	_, _ = sb.WriteString("some data")
	sb.Reset()
	require.Empty(t, sb.String())
	require.Equal(t, 0, sb.Len())
}

func TestSafeBuffer_Truncate(t *testing.T) {
	sb := NewSafeBuffer()
	_, _ = sb.WriteString("abcdef")
	sb.Truncate(3)
	require.Equal(t, "abc", sb.String())
}
