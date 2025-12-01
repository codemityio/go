package dsn

import "net/url"

// DSN represents a parsed data source name, supporting URL-based formats
// (e.g., mysql://, mongodb://, sqlite://) as well as file paths. It provides
// normalised access to scheme, credentials, host, port, path, query parameters,
// and optional cache size suffixes.
type DSN struct {
	Scheme   string
	User     string
	Password string
	Host     string   // first host, for convenience
	Port     int      // first port, for convenience
	Hosts    []string // all hosts, including the first
	Ports    []int    // all ports, including the first
	Path     string
	Query    url.Values
	Raw      *url.URL
}
