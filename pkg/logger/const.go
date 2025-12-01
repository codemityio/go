package logger

const (
	labelLevel   = "level"
	labelMessage = "message"
	labelError   = "error"

	// LevelDebug a debug level for the logger.
	LevelDebug Level = "debug"
	// LevelInfo an info level for the logger.
	LevelInfo Level = "info"
	// LevelWarn warning level for the logger.
	LevelWarn Level = "warning"
	// LevelError an error level for the logger.
	LevelError Level = "error"
	// LevelFatal an fatal level for the logger.
	LevelFatal Level = "fatal"

	// debug level for the logger.
	d uint = 1 << iota
	// info level for the logger.
	i
	// warning level for the logger.
	w
	// error level for the logger.
	e
	// fatal level for the logger.
	f

	// FormatText a simple text format for the log output.
	FormatText Format = "text"
	// FormatJSON a JSON format for the log output.
	FormatJSON Format = "json"

	// BufferSizeDefault async output buffer size.
	BufferSizeDefault = 1024
)
