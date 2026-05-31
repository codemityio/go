package generator

// New a fsm diagram generator factory.
func New(options ...Option) *DefaultGenerator {
	const black = "black"

	generator := &DefaultGenerator{
		start: black,
		stop:  black,
		state: "#f5f5f5",
		link:  black,
	}

	for _, opt := range options {
		opt(generator)
	}

	return generator
}
