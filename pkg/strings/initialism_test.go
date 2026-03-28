package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var list = []string{
	"acl",
	"api",
	"ascii",
	"cpu",
	"css",
	"dns",
	"eod",
	"eof",
	"eol",
	"guid",
	"html",
	"http",
	"https",
	"id",
	"ip",
	"json",
	"ok",
	"otp",
	"ram",
	"rpc",
	"sdk",
	"sms",
	"smtp",
	"sku",
	"sql",
	"ssh",
	"tcp",
	"tls",
	"ttl",
	"udp",
	"ui",
	"uid",
	"uuid",
	"uri",
	"url",
	"utf",
	"vm",
	"xml",
	"xmpp",
	"xsrf",
	"xss",
}

type testCase struct {
	input  string
	output string
}

func TestToPascal(t *testing.T) {
	tests := map[string]testCase{
		"hello": {
			input:  "hello",
			output: "Hello",
		},
		"hello-world": {
			input:  "hello-world",
			output: "HelloWorld",
		},
		"hello_world": {
			input:  "hello_world",
			output: "HelloWorld",
		},
		"HelloWorld": {
			input:  "HelloWorld",
			output: "HelloWorld",
		},
		"helloWorld": {
			input:  "helloWorld",
			output: "HelloWorld",
		},
		"hello-http-world": {
			input:  "hello-http-world",
			output: "HelloHTTPWorld",
		},
		"http-server-config": {
			input:  "http-server-config",
			output: "HTTPServerConfig",
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			initialism := NewInitialism(WithInitialismList(list))

			result := initialism.ToPascal(test.input)

			assert.Equal(t, test.output, result)
		})
	}
}

func TestToCamel(t *testing.T) {
	tests := map[string]testCase{
		"hello": {
			input:  "hello",
			output: "hello",
		},
		"hello-world": {
			input:  "hello-world",
			output: "helloWorld",
		},
		"hello_world": {
			input:  "hello_world",
			output: "helloWorld",
		},
		"HelloWorld": {
			input:  "HelloWorld",
			output: "helloWorld",
		},
		"helloWorld": {
			input:  "helloWorld",
			output: "helloWorld",
		},
		"http-server-config": {
			input:  "http-server-config",
			output: "httpServerConfig",
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			initialism := NewInitialism(WithInitialismList(list))

			result := initialism.ToCamel(test.input)

			assert.Equal(t, test.output, result)
		})
	}
}

func TestToKebab(t *testing.T) {
	tests := map[string]testCase{
		"Hello": {
			input:  "Hello",
			output: "hello",
		},
		"HelloWorld": {
			input:  "HelloWorld",
			output: "hello-world",
		},
		"hello-world": {
			input:  "hello-world",
			output: "hello-world",
		},
		"hello_world": {
			input:  "hello_world",
			output: "hello-world",
		},
		"helloWorld": {
			input:  "helloWorld",
			output: "hello-world",
		},
		"HTTPServerConfig": {
			input:  "HTTPServerConfig",
			output: "http-server-config",
		},
		"HTTPSServerConfig": {
			input:  "HTTPSServerConfig",
			output: "https-server-config",
		},
		"ServerHTTP": {
			input:  "ServerHTTP",
			output: "server-http",
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			initialism := NewInitialism(WithInitialismList(list))

			result := initialism.ToKebab(test.input)

			assert.Equal(t, test.output, result)
		})
	}
}

func TestToSnake(t *testing.T) {
	tests := map[string]testCase{
		"Hello": {
			input:  "Hello",
			output: "hello",
		},
		"HelloWorld": {
			input:  "HelloWorld",
			output: "hello_world",
		},
		"hello-world": {
			input:  "hello-world",
			output: "hello_world",
		},
		"hello_world": {
			input:  "hello_world",
			output: "hello_world",
		},
		"helloWorld": {
			input:  "helloWorld",
			output: "hello_world",
		},
		"HTTPServerConfig": {
			input:  "HTTPServerConfig",
			output: "http_server_config",
		},
		"HTTPSServerConfig": {
			input:  "HTTPSServerConfig",
			output: "https_server_config",
		},
		"ServerHTTP": {
			input:  "ServerHTTP",
			output: "server_http",
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			initialism := NewInitialism(WithInitialismList(list))

			result := initialism.ToSnake(test.input)

			assert.Equal(t, test.output, result)
		})
	}
}
