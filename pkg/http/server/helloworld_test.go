package server

import "time"

// GetHelloWorldTextPlainInput represents the input headers for the plain-text hello-world endpoint.
type GetHelloWorldTextPlainInput struct {
	Accept string `header:"Accept" validate:"required"`
}

// GetHelloWorldTextPlainOutput represents the plain-text response, including headers and body.
type GetHelloWorldTextPlainOutput struct {
	ContentType  string    `header:"Content-Type"  validate:"required"`
	LastModified time.Time `header:"Last-Modified"`
	Body         string    `                                           required:"true"`
}

// GetHelloWorldTextHTMLInput represents the input headers for the HTML hello-world endpoint.
type GetHelloWorldTextHTMLInput struct {
	Accept string `header:"Accept" validate:"required"`
}

// GetHelloWorldTextHTMLOutput represents the HTML response, including headers and body.
type GetHelloWorldTextHTMLOutput struct {
	ContentType  string    `header:"Content-Type"  validate:"required"`
	LastModified time.Time `header:"Last-Modified"`
	Body         string    `                                           required:"true"`
}

// GetHelloWorldApplicationJSONInput represents the input headers for the JSON hello-world endpoint.
type GetHelloWorldApplicationJSONInput struct {
	Accept string `header:"Accept" validate:"required"`
}

// GetHelloWorldApplicationJSONOutputBody represents the JSON response body for the hello-world endpoint.
type GetHelloWorldApplicationJSONOutputBody struct {
	Message string `doc:"Greeting message" example:"Hello, world!" json:"message"`
}

// GetHelloWorldApplicationJSONOutput represents the JSON response, including headers and body.
type GetHelloWorldApplicationJSONOutput struct {
	ContentType  string                                 `header:"Content-Type"  validate:"required"`
	LastModified time.Time                              `header:"Last-Modified"`
	Body         GetHelloWorldApplicationJSONOutputBody `                                           required:"true"`
}
