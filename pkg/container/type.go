package container

// Option to be used to build with optional deps.
type Option func(p *DefaultContainer)

// ShutdownFunc shutdown function.
type ShutdownFunc func()

// ServiceRecord service record structure.
type ServiceRecord struct {
	Service  any
	Shutdown func(service any)
}

// LogFunc to allow log warnings and infos.
type LogFunc func(warn string)
