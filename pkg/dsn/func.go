package dsn

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var bareSchemeRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+-.]*$`)

// Parse interprets the given string as a DSN, handling URL-style formats,
// file URLs, bare schemes, and relative or absolute file paths.
func Parse(input string) (*DSN, error) {
	input = strings.TrimSpace(input)

	dsn := &DSN{
		Scheme:   "",
		User:     "",
		Password: "",
		Host:     "",
		Port:     0,
		Hosts:    []string{},
		Ports:    []int{},
		Path:     "",
		Query:    url.Values{},
		Raw:      nil,
	}

	// 1. URL-based DSN (contains "://")
	if strings.Contains(input, "://") {
		// bare scheme like "memory://"
		if before, ok := strings.CutSuffix(input, "://"); ok {
			dsn.Scheme = before

			return dsn, nil
		}

		return parseURLDSN(input, dsn)
	}

	// 2. File paths (absolute or relative)
	if strings.Contains(input, "/") {
		return parseFilePath(input, dsn)
	}

	// 3. Bare scheme (memory, sqlite, etc.)
	if bareSchemeRe.MatchString(input) {
		dsn.Scheme = input

		return dsn, nil
	}

	return parseFilePath(input, dsn)
}

func parseFilePath(input string, dsn *DSN) (*DSN, error) {
	dsn.Scheme = "file"

	if before, after, ok := strings.Cut(input, "?"); ok {
		query, err := url.ParseQuery(after)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidQuery, err)
		}

		dsn.Path = before
		dsn.Query = query

		return dsn, nil
	}

	dsn.Path = input

	return dsn, nil
}

func parseURLDSN(input string, dsn *DSN) (*DSN, error) {
	// Normalise multi-host to single host before passing to url.Parse,
	// preserving the raw multi-host string for manual parsing below.
	normalised, rawHosts := normaliseMultiHost(input)

	parsed, err := url.Parse(normalised)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidDSN, err)
	}

	dsn.Raw = parsed
	dsn.Scheme = parsed.Scheme
	dsn.Query = parsed.Query()

	// Credentials
	if parsed.User != nil {
		dsn.User = parsed.User.Username()
		dsn.Password, _ = parsed.User.Password()
	}

	// Parse all hosts and ports
	hosts, ports, err := parseHosts(rawHosts)
	if err != nil {
		return nil, err
	}

	dsn.Hosts = hosts
	dsn.Ports = ports

	// First host/port for convenience
	if len(hosts) > 0 {
		dsn.Host = hosts[0]
		dsn.Port = ports[0]
	}

	dsn.Path = parsed.Path

	return dsn, nil
}

// normaliseMultiHost extracts the raw host section and returns a normalised
// single-host URL safe for url.Parse, along with the raw hosts string.
func normaliseMultiHost(input string) (string, string) {
	schemeEnd := strings.Index(input, "://")
	if schemeEnd == -1 {
		return input, ""
	}

	rest := input[schemeEnd+3:]

	slashIdx := strings.Index(rest, "/")

	authority := rest
	pathAndQuery := ""

	if slashIdx != -1 {
		authority = rest[:slashIdx]
		pathAndQuery = rest[slashIdx:]
	}

	userInfo := ""
	hosts := authority

	if atIdx := strings.Index(authority, "@"); atIdx != -1 {
		userInfo = authority[:atIdx+1]
		hosts = authority[atIdx+1:]
	}

	// Use only the first host for url.Parse
	firstHost := hosts
	if before, _, ok := strings.Cut(hosts, ","); ok {
		firstHost = before
	}

	return input[:schemeEnd+3] + userInfo + firstHost + pathAndQuery, hosts
}

// parseHosts splits a comma-separated host:port string into separate slices.
func parseHosts(rawHosts string) ([]string, []int, error) {
	hosts := make([]string, 0)
	ports := make([]int, 0)

	if rawHosts == "" {
		return hosts, ports, nil
	}

	for entry := range strings.SplitSeq(rawHosts, ",") {
		entry = strings.TrimSpace(entry)

		host := entry
		port := 0

		var err error

		if idx := strings.LastIndex(entry, ":"); idx != -1 {
			host = entry[:idx]

			port, err = strconv.Atoi(entry[idx+1:])
			if err != nil {
				return nil, nil, fmt.Errorf("%w: in %q: %w", ErrInvalidPort, entry, err)
			}
		}

		hosts = append(hosts, host)
		ports = append(ports, port)
	}

	return hosts, ports, nil
}
