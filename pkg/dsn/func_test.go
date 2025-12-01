package dsn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantScheme  string
		wantUser    string
		wantPass    string
		wantHost    string
		wantPort    int
		wantHosts   []string
		wantPorts   []int
		wantPath    string
		wantQuery   map[string]string
		expectError bool
	}{
		{
			name:       "bare scheme",
			input:      "memory",
			wantScheme: "memory",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "bare scheme with slashes",
			input:      "memory://",
			wantScheme: "memory",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "relative file path",
			input:      "configs/data.json",
			wantScheme: "file",
			wantPath:   "configs/data.json",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "absolute file path",
			input:      "/etc/app.json",
			wantScheme: "file",
			wantPath:   "/etc/app.json",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "parent relative file path",
			input:      "../data.json",
			wantScheme: "file",
			wantPath:   "../data.json",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "same directory path",
			input:      "./data.json",
			wantScheme: "file",
			wantPath:   "./data.json",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "file path with query",
			input:      "/var/broker?cacheSize=1024&mode=ro",
			wantScheme: "file",
			wantHost:   "",
			wantPort:   0,
			wantPath:   "/var/broker",
			wantHosts:  []string{},
			wantPorts:  []int{},
			wantQuery: map[string]string{
				"cacheSize": "1024",
				"mode":      "ro",
			},
		},
		{
			name:       "relative file path with query",
			input:      "var/broker?cacheSize=1024&mode=ro",
			wantScheme: "file",
			wantHost:   "",
			wantPort:   0,
			wantPath:   "var/broker",
			wantHosts:  []string{},
			wantPorts:  []int{},
			wantQuery: map[string]string{
				"cacheSize": "1024",
				"mode":      "ro",
			},
		},
		{
			name:       "file url absolute",
			input:      "file:///etc/app.json",
			wantScheme: "file",
			wantHost:   "",
			wantPort:   0,
			wantPath:   "/etc/app.json",
			wantHosts:  []string{},
			wantPorts:  []int{},
		},
		{
			name:       "file url with host",
			input:      "file://host/share/file.json",
			wantScheme: "file",
			wantHost:   "host",
			wantPort:   0,
			wantHosts:  []string{"host"},
			wantPorts:  []int{0},
			wantPath:   "/share/file.json",
		},
		{
			name:       "mysql with query params",
			input:      "mysql://user:pass@localhost:3306/mydb?timeout=5s&tls=skip",
			wantScheme: "mysql",
			wantUser:   "user",
			wantPass:   "pass",
			wantHost:   "localhost",
			wantPort:   3306,
			wantHosts:  []string{"localhost"},
			wantPorts:  []int{3306},
			wantPath:   "/mydb",
			wantQuery: map[string]string{
				"timeout": "5s",
				"tls":     "skip",
			},
		},
		{
			name:       "mysql with user pass and port",
			input:      "mysql://user:pass@localhost:3306/mydb",
			wantScheme: "mysql",
			wantUser:   "user",
			wantPass:   "pass",
			wantHost:   "localhost",
			wantPort:   3306,
			wantHosts:  []string{"localhost"},
			wantPorts:  []int{3306},
			wantPath:   "/mydb",
		},
		{
			name:       "mysql without password",
			input:      "mysql://user@localhost:3306/mydb",
			wantScheme: "mysql",
			wantUser:   "user",
			wantHost:   "localhost",
			wantPort:   3306,
			wantHosts:  []string{"localhost"},
			wantPorts:  []int{3306},
			wantPath:   "/mydb",
		},
		{
			name:       "mongodb single host",
			input:      "mongodb://root:secret@mongo.example.com:27017/appdb",
			wantScheme: "mongodb",
			wantUser:   "root",
			wantPass:   "secret",
			wantHost:   "mongo.example.com",
			wantPort:   27017,
			wantHosts:  []string{"mongo.example.com"},
			wantPorts:  []int{27017},
			wantPath:   "/appdb",
		},
		{
			name:       "mongodb multiple hosts",
			input:      "mongodb://root:secret@db1:27017,db2:27018/app?replicaSet=rs0",
			wantScheme: "mongodb",
			wantUser:   "root",
			wantPass:   "secret",
			wantHost:   "db1",
			wantPort:   27017,
			wantHosts:  []string{"db1", "db2"},
			wantPorts:  []int{27017, 27018},
			wantPath:   "/app",
			wantQuery: map[string]string{
				"replicaSet": "rs0",
			},
		},
		{
			name:        "invalid url",
			input:       "://bad",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Parse(tt.input)

			if tt.expectError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, out)

			require.Equal(t, tt.wantScheme, out.Scheme)
			require.Equal(t, tt.wantUser, out.User)
			require.Equal(t, tt.wantPass, out.Password)
			require.Equal(t, tt.wantHost, out.Host)
			require.Equal(t, tt.wantPort, out.Port)
			require.Equal(t, tt.wantPath, out.Path)
			require.Equal(t, tt.wantHosts, out.Hosts)
			require.Equal(t, tt.wantPorts, out.Ports)

			if len(tt.wantQuery) == 0 {
				require.Empty(t, out.Query)
			} else {
				for k, v := range tt.wantQuery {
					require.Equal(t, v, out.Query.Get(k), "query param %q mismatch", k)
				}
			}
		})
	}
}
