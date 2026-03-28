package logger

type (
	// Level a level type.
	Level string
	// Format a format type.
	Format string
	// Option to be used to build with optional deps.
	Option func(logger *DefaultLogger)

	logEntry struct {
		level   uint
		label   Level
		message string
		err     error
		fields  []Field
	}
)
