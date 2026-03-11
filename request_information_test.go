package cute

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/require"
	"net/url"
)

func TestRequestInformation(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://example.com/test"),
			},
		},
		RequestInformation: []RequestInformation{
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("custom_key", "custom_value"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	require.NotEmpty(t, capturedT.captured.params)
	require.Len(t, capturedT.captured.params, 1)
	require.Equal(t, "custom_value", getParameterValue(capturedT.captured.params, "custom_key"))
}

func TestRequestInformation_Error(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://example.com/test"),
			},
		},
		RequestInformation: []RequestInformation{
			func(req *http.Request) ([]*allure.Parameter, error) {
				return nil, errors.New("test error")
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	// Errors in RequestInformation are logged but do not fail the test
	// The request is still executed successfully
	require.Empty(t, results.GetErrors())
	// No parameters should be captureParams when the callback returns an error
	require.Empty(t, capturedT.captured.params)
}

func TestRequestInformation_MultipleCallbacks(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://example.com/test"),
			},
		},
		RequestInformation: []RequestInformation{
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("key1", "value1"), nil
			},
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("key2", "value2"), nil
			},
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("key3", "value3", "key4", "value4"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Should capture parameters from all callbacks: key1, key2, method, host
	require.Len(t, capturedT.captured.params, 4)
	require.Equal(t, "value1", getParameterValue(capturedT.captured.params, "key1"))
	require.Equal(t, "value2", getParameterValue(capturedT.captured.params, "key2"))
	require.Equal(t, "value3", getParameterValue(capturedT.captured.params, "key3"))
	require.Equal(t, "value4", getParameterValue(capturedT.captured.params, "key4"))
}

func TestRequestInformationT(t *testing.T) {
	t.Parallel()
	var capturedRequest *http.Request

	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodPost),
				WithURI("http://example.com/test"),
				WithBody([]byte("test body")),
			},
		},
		RequestInformationT: []RequestInformationT{
			func(t T, req *http.Request) ([]*allure.Parameter, error) {
				capturedRequest = req
				return allure.NewParameters("t_callback", "works"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	require.NotNil(t, capturedRequest)
	require.Equal(t, http.MethodPost, capturedRequest.Method)
	// Verify that the parameter from RequestInformationT was captureParams
	require.NotEmpty(t, capturedT.captured.params)
	require.Len(t, capturedT.captured.params, 1)
	require.Equal(t, "works", getParameterValue(capturedT.captured.params, "t_callback"))
}

func TestRequestInformationT_Error(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://example.com/test"),
			},
		},
		RequestInformationT: []RequestInformationT{
			func(t T, req *http.Request) ([]*allure.Parameter, error) {
				return nil, errors.New("t_callback error")
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	// Errors in RequestInformation are logged but do not fail the test
	// The request is still executed successfully
	require.Empty(t, results.GetErrors())
	// No parameters should be captureParams when the callback returns an error
	require.Empty(t, capturedT.captured.params)
}

func TestRequestInformationT_MultipleCallbacks(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://example.com/test"),
			},
		},
		RequestInformationT: []RequestInformationT{
			func(t T, req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("t_key1", "t_value1"), nil
			},
			func(t T, req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("t_key2", "t_value2"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Should capture parameters from both callbacks
	require.Len(t, capturedT.captured.params, 2)
	require.Equal(t, "t_value1", getParameterValue(capturedT.captured.params, "t_key1"))
	require.Equal(t, "t_value2", getParameterValue(capturedT.captured.params, "t_key2"))
}

func TestRequestInformationAndRequestInformationT_Together(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://example.com/test"),
			},
		},
		RequestInformation: []RequestInformation{
			func(req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("req_key", "req_value"), nil
			},
		},
		RequestInformationT: []RequestInformationT{
			func(t T, req *http.Request) ([]*allure.Parameter, error) {
				return allure.NewParameters("req_t_key", "req_t_value"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Should capture parameters from both RequestInformation and RequestInformationT
	require.Len(t, capturedT.captured.params, 2)
	require.Equal(t, "req_value", getParameterValue(capturedT.captured.params, "req_key"))
	require.Equal(t, "req_t_value", getParameterValue(capturedT.captured.params, "req_t_key"))
}

func TestWithNewStep_CapturesParameters(t *testing.T) {
	t.Parallel()
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Call WithNewStep with parameters that should be captured
	capturedT.WithNewStep("Test step", func(stepCtx provider.StepCtx) {
		// Parameters added within the step should also be captured
		stepCtx.WithParameters(allure.NewParameters("step_param", "step_value")...)
	}, allure.NewParameters("step_key", "step_value")...)

	// Verify both parameters passed to WithNewStep and added within it are captured
	require.NotEmpty(t, capturedT.captured.params)
	require.Len(t, capturedT.captured.params, 2)
	require.Equal(t, "step_value", getParameterValue(capturedT.captured.params, "step_key"))
	require.Equal(t, "step_value", getParameterValue(capturedT.captured.params, "step_param"))
}

// paramCapture holds captured parameters - shared across parent and child captureT instances
type paramCapture struct {
	params []*allure.Parameter
}

// captureT embeds provider.T and captures all parameters passed to WithParameters
type captureT struct {
	provider.T
	captured *paramCapture
}

// newCaptureT creates a new captureT wrapper
func newCaptureT(t provider.T) *captureT {
	return &captureT{
		T: t,
		captured: &paramCapture{
			params: make([]*allure.Parameter, 0),
		},
	}
}

// captureStepCtx wraps provider.StepCtx and captures parameters added within steps
type captureStepCtx struct {
	provider.StepCtx
	captured *paramCapture
}

// WithParameters captures parameters and forwards to embedded provider.StepCtx
func (cs *captureStepCtx) WithParameters(parameters ...*allure.Parameter) {
	cs.captured.params = append(cs.captured.params, parameters...)
	cs.StepCtx.WithParameters(parameters...)
}

// WithNewParameters forwards to embedded provider.StepCtx
func (cs *captureStepCtx) WithNewParameters(kv ...interface{}) {
	cs.StepCtx.WithNewParameters(kv...)
}

// WithParameters captures parameters and forwards to embedded provider.T
func (c *captureT) WithParameters(parameters ...*allure.Parameter) {
	c.captured.params = append(c.captured.params, parameters...)
	c.T.WithParameters(parameters...)
}

// WithNewParameters forwards to embedded provider.T
func (c *captureT) WithNewParameters(kv ...interface{}) {
	c.T.WithNewParameters(kv...)
}

// WithNewStep captures parameters and forwards to embedded provider.T
func (c *captureT) WithNewStep(name string, step func(sCtx provider.StepCtx), params ...*allure.Parameter) {
	// Capture parameters passed to WithNewStep
	c.captured.params = append(c.captured.params, params...)

	// Wrap the step callback to capture parameters added within the step
	wrappedStep := func(stepCtx provider.StepCtx) {
		wrappedCtx := &captureStepCtx{
			StepCtx:  stepCtx,
			captured: c.captured,
		}
		step(wrappedCtx)
	}

	c.T.WithNewStep(name, wrappedStep, params...)
}

// Run wraps the test body to capture nested parameters
func (c *captureT) Run(testName string, testBody func(provider.T), tags ...string) *allure.Result {
	return c.T.Run(testName, func(innerT provider.T) {
		// Wrap the inner T to capture its parameters too
		// Share the same paramCapture pointer so nested params are propagated back
		wrappedInnerT := &captureT{
			T:        innerT,
			captured: c.captured,
		}
		testBody(wrappedInnerT)
	}, tags...)
}

// Helper function to find parameter by key
func getParameterValue(params []*allure.Parameter, key string) any {
	for _, p := range params {
		if p.Name == key {
			return p.Value
		}
	}
	return nil
}

func TestRequestCurlInformation_ValidRequest(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodPost, "http://example.com/api", nil)
	require.NoError(t, err)

	params, err := RequestCurlInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)
	require.NotEmpty(t, params)

	curlParam := getParameterValue(params, "curl")
	require.NotNil(t, curlParam)
	require.NotEmpty(t, curlParam)
}

func TestRequestCurlInformation_InvalidRequest(t *testing.T) {
	t.Parallel()
	// Create request with invalid URL structure to trigger curl conversion error
	req := &http.Request{
		Method: http.MethodGet,
		URL:    nil, // Invalid: nil URL
	}

	params, err := RequestCurlInformation(req)

	// Should return error
	require.Error(t, err)
	require.Nil(t, params)
}

func TestRequestHTTPBaseInformation_StandardRequest(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodPost, "http://example.com/test", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBaseInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)
	require.Equal(t, 2, len(params))

	require.Equal(t, http.MethodPost, getParameterValue(params, "method"))
	require.Equal(t, "example.com", getParameterValue(params, "host"))
}

func TestRequestHTTPBaseInformation_DifferentMethods(t *testing.T) {
	t.Parallel()
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		req, err := http.NewRequest(method, "http://api.test.com/resource", nil)
		require.NoError(t, err)

		params, err := RequestHTTPBaseInformation(req)

		require.NoError(t, err)
		require.NotNil(t, params)
		require.Equal(t, method, getParameterValue(params, "method"))
	}
}

func TestRequestHTTPBaseInformation_EmptyHost(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Method: http.MethodGet,
		Host:   "",
		URL:    &url.URL{},
	}

	params, err := RequestHTTPBaseInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)
	require.Equal(t, "", getParameterValue(params, "host"))
}

func TestRequestHTTPHeadersInformation_WithHeaders(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")

	params, err := RequestHTTPHeadersInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	headersValue := getParameterValue(params, "headers")
	require.NotNil(t, headersValue)
	require.NotEmpty(t, headersValue)
	// Should contain headers as JSON string
	headersStr := headersValue.(string)
	require.Contains(t, headersStr, "Content-Type")
	require.Contains(t, headersStr, "Authorization")
}

func TestRequestHTTPHeadersInformation_NoHeaders(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	params, err := RequestHTTPHeadersInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	headersValue := getParameterValue(params, "headers")
	require.NotNil(t, headersValue)
}

func TestRequestHTTPHeadersInformation_MultipleHeaderValues(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "text/html")

	params, err := RequestHTTPHeadersInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	headersValue := getParameterValue(params, "headers")
	require.NotNil(t, headersValue)
	require.NotEmpty(t, headersValue)
}

func TestRequestHTTPBodyInformation_NoBody(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBodyInformation(req)

	require.NoError(t, err)
	// Should return nil when Body is nil
	require.Nil(t, params)
}

func TestRequestHTTPBodyInformation_EmptyBody(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBodyInformation(req)

	require.NoError(t, err)
	// Empty/nil body should return nil
	require.Nil(t, params)
}

func TestRequestHTTPBaseInformation_PortInHost(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://example.com:8080/test", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBaseInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	hostValue := getParameterValue(params, "host")
	require.NotNil(t, hostValue)
	// Host should include port when present in URL
	require.Contains(t, hostValue.(string), "example.com")
}

func TestRequestHTTPBaseInformation_SubdomainRequest(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://api.v1.example.com/endpoint", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBaseInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	hostValue := getParameterValue(params, "host")
	require.Equal(t, "api.v1.example.com", hostValue)
}

func TestRequestHTTPBaseInformation_HttpsRequest(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "https://secure.example.com/endpoint", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBaseInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	methodValue := getParameterValue(params, "method")
	hostValue := getParameterValue(params, "host")
	require.NotNil(t, methodValue)
	require.NotNil(t, hostValue)
	// Scheme should not affect the returned parameters
	require.Equal(t, "secure.example.com", hostValue)
}

// Tests for sanitized request data in RequestInformation and RequestInformationT

func TestRequestInformation_ReceivesSanitizedData(t *testing.T) {
	t.Parallel()

	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("https://example.com/api"),
				WithHeaders(map[string][]string{
					"Authorization": {"Bearer secret_token_abc123"},
				}),
			},
		},
		RequestInformation: []RequestInformation{
			func(req *http.Request) ([]*allure.Parameter, error) {
				// This callback receives the sanitized request
				authHeader := req.Header.Get("Authorization")
				return allure.NewParameters("auth", authHeader), nil
			},
		},
		// Sanitizer replaces sensitive header before RequestInformation is called
		RequestSanitizer: func(req *http.Request) {
			req.Header.Set("Authorization", "Bearer [REDACTED]")
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())

	require.Equal(t, "Bearer [REDACTED]", getParameterValue(capturedT.captured.params, "auth"))
}

func TestRequestInformationT_ReceivesSanitizedData(t *testing.T) {
	t.Parallel()

	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodPost),
				WithURI("https://example.com/login"),
				WithHeaders(map[string][]string{
					"Authorization": {"Bearer secret_token_abc123"},
				}),
			},
		},
		RequestInformationT: []RequestInformationT{
			func(t T, req *http.Request) ([]*allure.Parameter, error) {
				// This callback receives the sanitized request
				authHeader := req.Header.Get("Authorization")
				return allure.NewParameters("auth", authHeader), nil
			},
		},
		// Sanitizer replaces sensitive header before RequestInformationT is called
		RequestSanitizer: func(req *http.Request) {
			req.Header.Set("Authorization", "Bearer [REDACTED]")
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())

	require.Equal(t, "Bearer [REDACTED]", getParameterValue(capturedT.captured.params, "auth"))
}
