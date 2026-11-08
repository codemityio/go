package validator

// New validator factory function.
func New() *DefaultValidator {
	return &DefaultValidator{}
}

// NewGeneric generic validator factory function.
func NewGeneric[T any]() *GenericValidator[T] {
	return &GenericValidator[T]{
		base: DefaultValidator{},
	}
}
