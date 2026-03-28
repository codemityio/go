package generator

import "strings"

// WithColours configuration option.
func WithColours(start, stop, state, link string) Option {
	return func(w *DefaultGenerator) {
		w.start = strings.Trim(start, "#")
		w.stop = strings.Trim(stop, "#")
		w.state = strings.Trim(state, "#")
		w.link = strings.Trim(link, "#")
	}
}
