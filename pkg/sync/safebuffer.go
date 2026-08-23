package sync

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

// SafeBuffer is a thread-safe wrapper around bytes.Buffer.
type SafeBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

// IsSafeBufferWrite reports whether err is or wraps ErrSafeBufferWrite.
func (s *SafeBuffer) IsSafeBufferWrite(err error) bool {
	return errors.Is(err, ErrSafeBufferWrite)
}

// Write writes to the buffer.
func (s *SafeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.b.Write(p)
	if err != nil {
		return 0, fmt.Errorf("%w:%w", ErrSafeBufferWrite, err)
	}

	return res, nil
}

// String returns the buffer content as a string.
func (s *SafeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.b.String()
}

// Bytes returns a copy of the buffer contents.
func (s *SafeBuffer) Bytes() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]byte(nil), s.b.Bytes()...)
}

// Len returns the number of bytes in the buffer.
func (s *SafeBuffer) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.b.Len()
}

// Cap returns the capacity of the buffer.
func (s *SafeBuffer) Cap() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.b.Cap()
}

// Reset resets the buffer to be empty.
func (s *SafeBuffer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.b.Reset()
}

// WriteString writes a string to the buffer.
func (s *SafeBuffer) WriteString(str string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.b.WriteString(str)
	if err != nil {
		return 0, fmt.Errorf("%w:%w", ErrSafeBufferWrite, err)
	}

	return res, nil
}

// WriteRune writes a single UTF-8 encoded rune to the buffer.
func (s *SafeBuffer) WriteRune(r rune) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.b.WriteRune(r)
	if err != nil {
		return 0, fmt.Errorf("%w:%w", ErrSafeBufferWrite, err)
	}

	return res, nil
}

// WriteByte writes a single byte to the buffer.
func (s *SafeBuffer) WriteByte(c byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if e := s.b.WriteByte(c); e != nil {
		return fmt.Errorf("%w:%w", ErrSafeBufferWrite, e)
	}

	return nil
}

// Truncate truncates the buffer to the specified length.
func (s *SafeBuffer) Truncate(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.b.Truncate(n)
}
