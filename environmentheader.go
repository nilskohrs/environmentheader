// Package environmentheader a plugin to use environment variables in headers
package environmentheader

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
	RequestHeaders  []HeaderMapping `json:"requestHeaders,omitempty"`
	ResponseHeaders []HeaderMapping `json:"responseHeaders,omitempty"`
}

// HeaderMapping is part of the plugin configuration.
type HeaderMapping struct {
	Header   string `json:"header,omitempty"`
	Env      string `json:"env,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

type environmentHeaderPlugin struct {
	RequestHeaders  []HeaderMapping
	ResponseHeaders []HeaderMapping
	next            http.Handler
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates a new EnvironmentHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	requestHeaders := make([]HeaderMapping, 0, len(config.RequestHeaders))
	for _, requestHeader := range config.RequestHeaders {
		if len(requestHeader.Header) == 0 {
			return nil, fmt.Errorf("missing header parameter on a request header mapping")
		}
		if len(requestHeader.Env) == 0 {
			return nil, fmt.Errorf("missing env parameter for request header `%s`", requestHeader.Header)
		}
		environmentVar := requestHeader.Env
		requestHeader.Env = os.Getenv(environmentVar)
		if !requestHeader.Optional && len(requestHeader.Env) == 0 {
			return nil, fmt.Errorf("environment variable `%s` is not set for request header `%s`", environmentVar, requestHeader.Header)
		}
		requestHeaders = append(requestHeaders, requestHeader)
	}
	responseHeaders := make([]HeaderMapping, 0, len(config.ResponseHeaders))
	for _, responseHeader := range config.ResponseHeaders {
		if len(responseHeader.Header) == 0 {
			return nil, fmt.Errorf("missing header parameter on a response header mapping")
		}
		if len(responseHeader.Env) == 0 {
			return nil, fmt.Errorf("missing env parameter for response header `%s`", responseHeader.Header)
		}
		environmentVar := responseHeader.Env
		responseHeader.Env = os.Getenv(environmentVar)
		if !responseHeader.Optional && len(responseHeader.Env) == 0 {
			return nil, fmt.Errorf("environment variable `%s` is not set for response header `%s`", environmentVar, responseHeader.Header)
		}
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
	for _, responseHeader := range c.ResponseHeaders {
		rw.Header().Add(responseHeader.Header, responseHeader.Env)
	}
	c.next.ServeHTTP(rw, req)
}
