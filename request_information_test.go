package cute

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/stretchr/testify/require"
)

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
func findParameter(params []*allure.Parameter, key string) *allure.Parameter {
	for _, p := range params {
		if p.Name == key {
			return p
		}
	}
	return nil
}

func TestRequestInformation(t *testing.T) {
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
	require.Equal(t, 1, len(capturedT.captured.params))
	require.Equal(t, "custom_key", capturedT.captured.params[0].Name)
	require.Equal(t, "custom_value", capturedT.captured.params[0].Value)
}

func TestRequestInformation_Error(t *testing.T) {
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
			RequestHTTPBaseInformation,
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Should capture parameters from all callbacks: key1, key2, method, host
	require.GreaterOrEqual(t, len(capturedT.captured.params), 4)
	require.NotNil(t, findParameter(capturedT.captured.params, "key1"))
	require.NotNil(t, findParameter(capturedT.captured.params, "key2"))
	require.NotNil(t, findParameter(capturedT.captured.params, "method"))
	require.NotNil(t, findParameter(capturedT.captured.params, "host"))
}

func TestRequestInformationT(t *testing.T) {
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
	require.Equal(t, 1, len(capturedT.captured.params))
	require.Equal(t, "t_callback", capturedT.captured.params[0].Name)
	require.Equal(t, "works", capturedT.captured.params[0].Value)
}

func TestRequestInformationT_Error(t *testing.T) {
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
	require.Equal(t, 2, len(capturedT.captured.params))
	require.Equal(t, "t_key1", capturedT.captured.params[0].Name)
	require.Equal(t, "t_value1", capturedT.captured.params[0].Value)
	require.Equal(t, "t_key2", capturedT.captured.params[1].Name)
	require.Equal(t, "t_value2", capturedT.captured.params[1].Value)
}

func TestRequestInformationAndRequestInformationT_Together(t *testing.T) {
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
	require.Equal(t, 2, len(capturedT.captured.params))
	require.NotNil(t, findParameter(capturedT.captured.params, "req_key"))
	require.NotNil(t, findParameter(capturedT.captured.params, "req_t_key"))
}

func TestWithNewStep_CapturesParameters(t *testing.T) {
	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// Call WithNewStep with parameters that should be captured
	capturedT.WithNewStep("Test step", func(stepCtx provider.StepCtx) {
		// Parameters added within the step should also be captured
		stepCtx.WithParameters(allure.NewParameters("step_param", "step_value")...)
	}, allure.NewParameters("step_key", "step_value")...)

	// Verify both parameters passed to WithNewStep and added within it are captured
	require.NotEmpty(t, capturedT.captured.params)
	require.GreaterOrEqual(t, len(capturedT.captured.params), 2)
	require.NotNil(t, findParameter(capturedT.captured.params, "step_key"))
	require.NotNil(t, findParameter(capturedT.captured.params, "step_param"))
}
