package cute

import (
	"net/http"
	"time"
)

type options struct {
	httpClient       *http.Client
	httpTimeout      time.Duration
	httpRoundTripper http.RoundTripper

	jsonMarshaler JSONMarshaler

	middleware *Middleware

	requestInformation   []RequestInformation
	requestInformationT  []RequestInformationT
	responseInformation  []ResponseInformation
	responseInformationT []ResponseInformationT
}

// Option ...
type Option func(*options)

// WithHTTPClient is a function for set custom http client
func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}

// WithJSONMarshaler is a function for set custom json marshaler
func WithJSONMarshaler(m JSONMarshaler) Option {
	return func(o *options) {
		o.jsonMarshaler = m
	}
}

// WithCustomHTTPTimeout is a function for set custom http client timeout
func WithCustomHTTPTimeout(t time.Duration) Option {
	return func(o *options) {
		o.httpTimeout = t
	}
}

// WithCustomHTTPRoundTripper is a function for set custom http round tripper
func WithCustomHTTPRoundTripper(r http.RoundTripper) Option {
	return func(o *options) {
		o.httpRoundTripper = r
	}
}

// WithMiddlewareAfter is function for set function which will run AFTER test execution
func WithMiddlewareAfter(after ...AfterExecute) Option {
	return func(o *options) {
		o.middleware.After = append(o.middleware.After, after...)
	}
}

// WithMiddlewareAfterT is function for set function which will run AFTER test execution
func WithMiddlewareAfterT(after ...AfterExecuteT) Option {
	return func(o *options) {
		o.middleware.AfterT = append(o.middleware.AfterT, after...)
	}
}

// WithMiddlewareBefore is function for set function which will run BEFORE test execution
func WithMiddlewareBefore(before ...BeforeExecute) Option {
	return func(o *options) {
		o.middleware.Before = append(o.middleware.Before, before...)
	}
}

// WithMiddlewareBeforeT is function for set function which will run BEFORE test execution
func WithMiddlewareBeforeT(beforeT ...BeforeExecuteT) Option {
	return func(o *options) {
		o.middleware.BeforeT = append(o.middleware.BeforeT, beforeT...)
	}
}

// WithRequestInformation sets request information handlers at HTTPTestMaker level
func WithRequestInformation(handlers ...RequestInformation) Option {
	return func(o *options) {
		o.requestInformation = append(o.requestInformation, handlers...)
	}
}

// WithRequestInformationT sets request information handlers with T context
func WithRequestInformationT(handlers ...RequestInformationT) Option {
	return func(o *options) {
		o.requestInformationT = append(o.requestInformationT, handlers...)
	}
}

// WithResponseInformation sets response information handlers at HTTPTestMaker level
func WithResponseInformation(handlers ...ResponseInformation) Option {
	return func(o *options) {
		o.responseInformation = append(o.responseInformation, handlers...)
	}
}

// WithResponseInformationT sets response information handlers with T context
func WithResponseInformationT(handlers ...ResponseInformationT) Option {
	return func(o *options) {
		o.responseInformationT = append(o.responseInformationT, handlers...)
	}
}
