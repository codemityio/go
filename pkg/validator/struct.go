package validator

// Error represents a validation error found.
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
