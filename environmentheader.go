// Package environmentheader a plugin to use environment variables in headers
package environmentheader

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/http/httpguts"
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
	for i := range config.RequestHeaders {
		err := loadData(&config.RequestHeaders[i], "request")
		if err != nil {
			return nil, err
		}
		requestHeaders = append(requestHeaders, config.RequestHeaders[i])
	}
	responseHeaders := make([]HeaderMapping, 0, len(config.ResponseHeaders))
	for i := range config.ResponseHeaders {
		err := loadData(&config.ResponseHeaders[i], "response")
		if err != nil {
			return nil, err
		}
		responseHeaders = append(responseHeaders, config.ResponseHeaders[i])
	}
	return &environmentHeaderPlugin{
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
		next:            next,
	}, nil
}

func loadData(requestHeader *HeaderMapping, headerType string) error {
	if !httpguts.ValidHeaderFieldName(requestHeader.Header) {
		return fmt.Errorf("%s header `%s` is an invalid header name", headerType, requestHeader.Header)
	}
	if len(requestHeader.Env) == 0 {
		return fmt.Errorf("missing env parameter for %s header `%s`", headerType, requestHeader.Header)
	}
	environmentVar := requestHeader.Env
	requestHeader.Env = os.Getenv(environmentVar)
	if !requestHeader.Optional && len(requestHeader.Env) == 0 {
		return fmt.Errorf("environment variable `%s` is empty for %s header `%s`", environmentVar, headerType, requestHeader.Header)
	}
	if !httpguts.ValidHeaderFieldValue(requestHeader.Env) {
		return fmt.Errorf("environment variable `%s` for %s header `%s` had an value which is not allowed as a header field value", environmentVar, headerType, requestHeader.Header)
	}
	return nil
}

func (c *environmentHeaderPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println()

	for _, requestHeader := range c.RequestHeaders {
		req.Header.Add(requestHeader.Header, requestHeader.Env)
	}
	for _, responseHeader := range c.ResponseHeaders {
		rw.Header().Add(responseHeader.Header, responseHeader.Env)
	}
	c.next.ServeHTTP(rw, req)
}
