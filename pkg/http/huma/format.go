package huma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"
)

// FormatTextPlain implements the Huma Format interface for plain text content.
type FormatTextPlain struct{}

// ContentType returns the MIME type associated with this format.
func (f *FormatTextPlain) ContentType() string {
	return ContentTypeTextPlain
}

// Marshal writes the encoded form of v to writer.
func (f *FormatTextPlain) Marshal(writer io.Writer, v any) error {
	res, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMarshal, err)
	}

	if _, e := io.Writer.Write(writer, res); e != nil {
		return fmt.Errorf("%w: %w", ErrMarshal, e)
	}

	return nil
}

// Unmarshal decodes the provided data into v.
//
// Plain text cannot be generically unmarshaled into structured Go types,
// so this returns an error for all cases.
func (f *FormatTextPlain) Unmarshal(_ []byte, _ any) error {
	return fmt.Errorf("%w: plain text unmarshaling is not supported", ErrUnmarshal)
}

// FormatTextHTML implements the Huma Format interface for a specific content type.
type FormatTextHTML struct {
	template *template.Template
}

// ContentType returns the MIME type associated with this format.
func (f *FormatTextHTML) ContentType() string {
	return ContentTypeTextHTML
}

// Marshal writes the encoded form of v to writer.
func (f *FormatTextHTML) Marshal(writer io.Writer, value any) error {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("%w: %w", ErrMarshal, err)
	}

	// remove trailing newline added by Encoder
	jsonStr := strings.TrimSuffix(buf.String(), "\n")

	if err := f.template.Execute(writer, jsonStr); err != nil {
		return fmt.Errorf("%w: %w", ErrMarshal, err)
	}

	return nil
}

// Unmarshal decodes the provided data into v.
// HTML input is not typically parsed into structs, so this is a no-op.
func (f *FormatTextHTML) Unmarshal(_ []byte, _ any) error {
	return fmt.Errorf("%w: HTML is not supported", ErrUnmarshal)
}

// FormatApplicationJSON implements the Huma Format interface for a specific content type.
type FormatApplicationJSON struct{}

// ContentType returns the MIME type associated with this format.
func (f *FormatApplicationJSON) ContentType() string {
	return ContentTypeApplicationJSON
}

// Marshal writes the encoded form of v to writer.
func (f *FormatApplicationJSON) Marshal(writer io.Writer, v any) error {
	if e := json.NewEncoder(writer).Encode(v); e != nil {
		return fmt.Errorf("%w: %w", ErrMarshal, e)
	}

	return nil
}

// Unmarshal decodes the provided data into v.
func (f *FormatApplicationJSON) Unmarshal(data []byte, v any) error {
	if e := json.Unmarshal(data, v); e != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshal, e)
	}

	return nil
}
