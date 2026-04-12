package transform

// Transformer generic transform interface.
type Transformer interface {
	// Transform converts input into output. The output parameter must be a pointer.
	Transform(input any, output any) error
}
