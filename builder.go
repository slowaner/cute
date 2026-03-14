package cute

import (
	"net/http"
	"slices"
	"time"
)

const defaultHTTPTimeout = 30

var (
	errorAssertIsNil = "assert must be not nil"
)

// HTTPTestMaker is a creator tests
type HTTPTestMaker struct {
	httpClient    *http.Client
	middleware    *Middleware
	jsonMarshaler JSONMarshaler

	requestInformation   []RequestInformation
	requestInformationT  []RequestInformationT
	responseInformation  []ResponseInformation
	responseInformationT []ResponseInformationT
}

// NewHTTPTestMaker is function for set options for all cute.
// For example, you can set timeout for all requests or set custom http client.
// These options are applied globally to all tests created via NewTestBuilder().
//
// Options:
// - WithCustomHTTPTimeout - set timeout for all requests
// - WithHTTPClient - set custom http client
// - WithCustomHTTPRoundTripper - set custom http round tripper
// - WithJSONMarshaler - set custom json marshaler
// - WithMiddlewareAfter - set function which will run AFTER test execution
// - WithMiddlewareAfterT - set function which will run AFTER test execution with TB
// - WithMiddlewareBefore - set function which will run BEFORE test execution
// - WithMiddlewareBeforeT - set function which will run BEFORE test execution with TB
// - WithRequestInformation - add request information handlers to capture request details
// - WithRequestInformationT - add request information handlers with T context
// - WithResponseInformation - add response information handlers to capture response details
// - WithResponseInformationT - add response information handlers with T context
//
// Information handlers are executed for each test and their results are added to Allure reports.
// Handlers can also be added per-test using builder methods like RequestInformation() and ResponseInformation().
func NewHTTPTestMaker(opts ...Option) *HTTPTestMaker {
	var (
		o = &options{
			middleware: new(Middleware),
		}

		timeout                    = defaultHTTPTimeout * time.Second
		roundTripper               = http.DefaultTransport
		jsMarshaler  JSONMarshaler = &jsonMarshaler{}
	)

	for _, opt := range opts {
		opt(o)
	}

	if o.httpTimeout != 0 {
		timeout = o.httpTimeout
	}

	if o.httpRoundTripper != nil { //nolint
		roundTripper = o.httpRoundTripper
	}

	httpClient := &http.Client{
		Transport: roundTripper,
		Timeout:   timeout,
	}

	if o.httpClient != nil {
		httpClient = o.httpClient
	}

	if o.jsonMarshaler != nil {
		jsMarshaler = o.jsonMarshaler
	}

	m := &HTTPTestMaker{
		httpClient:           httpClient,
		jsonMarshaler:        jsMarshaler,
		middleware:           o.middleware,
		requestInformation:   o.requestInformation,
		requestInformationT:  o.requestInformationT,
		responseInformation:  o.responseInformation,
		responseInformationT: o.responseInformationT,
	}

	return m
}

// NewTestBuilder is a function for initialization foundation for cute
func (m *HTTPTestMaker) NewTestBuilder() AllureBuilder {
	tests := createDefaultTests(m)

	return &cute{
		baseProps:    m,
		countTests:   0,
		tests:        tests,
		allureInfo:   new(AllureInformation),
		allureLinks:  new(AllureLinks),
		allureLabels: new(AllureLabels),
		parallel:     false,
	}
}

func createInformationHandlersFromTemplate(
	m *HTTPTestMaker,
	additionalReqInfo []RequestInformation,
	additionalReqInfoT []RequestInformationT,
	additionalRespInfo []ResponseInformation,
	additionalRespInfoT []ResponseInformationT,
) (
	reqInfo []RequestInformation,
	reqInfoT []RequestInformationT,
	respInfo []ResponseInformation,
	respInfoT []ResponseInformationT,
) {
	// Deep copy request handlers using slices.Clone, ensure non-nil slices
	reqInfo = slices.Concat(m.requestInformation, additionalReqInfo)
	if reqInfo == nil {
		reqInfo = make([]RequestInformation, 0)
	}
	reqInfoT = slices.Concat(m.requestInformationT, additionalReqInfoT)
	if reqInfoT == nil {
		reqInfoT = make([]RequestInformationT, 0)
	}

	// Deep copy response handlers using slices.Clone, ensure non-nil slices
	respInfo = slices.Concat(m.responseInformation, additionalRespInfo)
	if respInfo == nil {
		respInfo = make([]ResponseInformation, 0)
	}
	respInfoT = slices.Concat(m.responseInformationT, additionalRespInfoT)
	if respInfoT == nil {
		respInfoT = make([]ResponseInformationT, 0)
	}

	return
}

func createDefaultTests(m *HTTPTestMaker) []*Test {
	tests := make([]*Test, 1)
	tests[0] = createDefaultTest(m, nil, nil, nil, nil)

	return tests
}

func createDefaultTest(
	m *HTTPTestMaker,
	reqInfo []RequestInformation,
	reqInfoT []RequestInformationT,
	respInfo []ResponseInformation,
	respInfoT []ResponseInformationT,
) *Test {
	totalReqInfo, totalReqInfoT, totalRespInfo, totalRespInfoT := createInformationHandlersFromTemplate(
		m,
		reqInfo,
		reqInfoT,
		respInfo,
		respInfoT,
	)

	return &Test{
		httpClient:    m.httpClient,
		jsonMarshaler: m.jsonMarshaler,
		Middleware:    createMiddlewareFromTemplate(m.middleware),
		AllureStep:    new(AllureStep),
		Request: &Request{
			Retry: new(RequestRetryPolitic),
		},
		Expect:               &Expect{JSONSchema: new(ExpectJSONSchema)},
		RequestInformation:   totalReqInfo,
		RequestInformationT:  totalReqInfoT,
		ResponseInformation:  totalRespInfo,
		ResponseInformationT: totalRespInfoT,
	}
}

func createMiddlewareFromTemplate(m *Middleware) *Middleware {
	after := make([]AfterExecute, 0, len(m.After))
	after = append(after, m.After...)

	afterT := make([]AfterExecuteT, 0, len(m.AfterT))
	afterT = append(afterT, m.AfterT...)

	before := make([]BeforeExecute, 0, len(m.Before))
	before = append(before, m.Before...)

	beforeT := make([]BeforeExecuteT, 0, len(m.BeforeT))
	beforeT = append(beforeT, m.BeforeT...)

	middleware := &Middleware{
		After:   after,
		AfterT:  afterT,
		Before:  before,
		BeforeT: beforeT,
	}

	return middleware
}

func (qt *cute) Create() MiddlewareRequest {
	return qt
}

func (qt *cute) CreateStep(name string) MiddlewareRequest {
	qt.tests[qt.countTests].AllureStep.Name = name

	return qt
}

func (qt *cute) CreateRequest() RequestHTTPBuilder {
	return qt
}
