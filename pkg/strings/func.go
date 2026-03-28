package strings

// WithInitialismList set initialism.
func WithInitialismList(list []string) InitialismOption {
	return func(l *Initialism) {
		l.list = make(map[string]struct{})

		for _, v := range list {
			l.list[v] = struct{}{}
		}
	}
}
