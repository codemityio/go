package server

// Option function.
type Option func(p *DefaultServer)

// LogFunc to allow log warnings and infos.
type LogFunc func(message string)
