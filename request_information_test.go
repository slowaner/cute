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
				WithURI(testServerAddress + "/test"),
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
				WithURI(testServerAddress + "/test"),
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
				WithURI(testServerAddress + "/test"),
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
				WithURI(testServerAddress + "/test"),
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
				WithURI(testServerAddress + "/test"),
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
				WithURI(testServerAddress + "/test"),
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
				WithURI(testServerAddress + "/test"),
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

func TestRequestCurlInformation_ValidRequest(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodPost, testServerAddress+"/api", nil)
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
	req, err := http.NewRequest(http.MethodPost, testServerAddress+"/test", nil)
	require.NoError(t, err)

	params, err := RequestHTTPBaseInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)
	require.Equal(t, 2, len(params))

	require.Equal(t, http.MethodPost, getParameterValue(params, "method"))
	require.Equal(t, testServerHost, getParameterValue(params, "host"))
}

func TestRequestHTTPBaseInformation_DifferentMethods(t *testing.T) {
	t.Parallel()
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, testServerAddress+"/resource", nil)
			require.NoError(t, err)

			params, err := RequestHTTPBaseInformation(req)

			require.NoError(t, err)
			require.NotNil(t, params)
			require.Equal(t, method, getParameterValue(params, "method"))
		})
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
	param := getParameterValue(params, "host")
	require.NotNil(t, param)
	require.Empty(t, param)
}

func TestRequestHTTPHeadersInformation_WithHeaders(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, testServerAddress, nil)
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
	req, err := http.NewRequest(http.MethodGet, testServerAddress, nil)
	require.NoError(t, err)

	params, err := RequestHTTPHeadersInformation(req)

	require.NoError(t, err)
	require.NotNil(t, params)

	headersValue := getParameterValue(params, "headers")
	require.NotNil(t, headersValue)
}

func TestRequestHTTPHeadersInformation_MultipleHeaderValues(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, testServerAddress, nil)
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
	req, err := http.NewRequest(http.MethodGet, testServerAddress, http.NoBody)
	require.NoError(t, err)

	params, err := RequestHTTPBodyInformation(req)

	require.NoError(t, err)
	// Should return nil when Body is nil
	require.Nil(t, params)
}

func TestRequestHTTPBodyInformation_NilBody(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodPost, testServerAddress, nil)
	require.NoError(t, err)

	params, err := RequestHTTPBodyInformation(req)

	require.NoError(t, err)
	// Empty/nil body should return nil
	require.Nil(t, params)
}

// Tests for sanitized request data in RequestInformation and RequestInformationT

func TestRequestInformation_ReceivesSanitizedData(t *testing.T) {
	t.Parallel()

	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/api"),
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
				WithURI(testServerAddress + "/login"),
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
