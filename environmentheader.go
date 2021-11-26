// Package environmentheader a plugin to use environment variables in headers
package environmentheader

import (
	"context"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
	RequestHeaders  []RequestHeader  `json:"requestHeaders,omitempty"`
	ResponseHeaders []ResponseHeader `json:"responseHeaders,omitempty"`
}

// RequestHeader is part of the plugin configuration.
type RequestHeader struct {
	Header string `json:"header,omitempty"`
	Env    string `json:"env,omitempty"`
}

// ResponseHeader is part of the plugin configuration.
type ResponseHeader struct {
	Header string `json:"header,omitempty"`
	Env    string `json:"env,omitempty"`
}

type environmentHeaderPlugin struct {
	RequestHeaders  []RequestHeader
	ResponseHeaders []ResponseHeader
	next            http.Handler
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates a new EnvironmentHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	requestHeaders := make([]RequestHeader, 0, len(config.RequestHeaders))
	for _, requestHeader := range config.RequestHeaders {
		requestHeader.Env = os.Getenv(requestHeader.Env)
		requestHeaders = append(requestHeaders, requestHeader)
	}
	responseHeaders := make([]ResponseHeader, 0, len(config.ResponseHeaders))
	for _, responseHeader := range config.ResponseHeaders {
		responseHeader.Env = os.Getenv(responseHeader.Env)
		responseHeaders = append(responseHeaders, responseHeader)
	}
	return &environmentHeaderPlugin{
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
		next:            next,
	}, nil
}

func (c *environmentHeaderPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, requestHeader := range c.RequestHeaders {
		req.Header.Add(requestHeader.Header, requestHeader.Env)
	}
	c.next.ServeHTTP(rw, req)
}
