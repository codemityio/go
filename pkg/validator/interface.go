package validator

//go:generate mockgen -source interface.go -destination interface_mock_test.go -package validator

// JSONValidator validator interface.
type JSONValidator interface {
	ValidateJSON(schema, input []byte) ([]Error, error)
}

// YAMLValidator validator interface.
type YAMLValidator interface {
	ValidateYAML(schema, input []byte) ([]Error, error)
}

// AnyValidator validator interface.
type AnyValidator interface {
	ValidateAny(schema []byte, input any) ([]Error, error)
}
