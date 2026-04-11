package client

import (
	"io"
	"net/http"
	"net/url"
)

// Clienter http client interface.
type Clienter interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
}
