package cute

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/require"
)

func TestResponseInformation(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/test"),
			},
		},
		ResponseInformation: []ResponseInformation{
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("response_custom_key", "response_custom_value"), nil
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
	require.Equal(t, "response_custom_value", getParameterValue(capturedT.captured.params, "response_custom_key"))
}

func TestResponseInformation_Error(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/test"),
			},
		},
		ResponseInformation: []ResponseInformation{
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return nil, errors.New("test error")
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	// Errors in ResponseInformation are logged but do not fail the test
	// The request is still executed successfully
	require.Empty(t, results.GetErrors())
	// No parameters should be captured when the callback returns an error
	require.Empty(t, capturedT.captured.params)
}

func TestResponseInformation_MultipleCallbacks(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/test"),
			},
		},
		ResponseInformation: []ResponseInformation{
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("resp_key1", "resp_value1"), nil
			},
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("resp_key2", "resp_value2"), nil
			},
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("resp_key3", "resp_value3", "resp_key4", "resp_value4"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Should capture parameters from all callbacks
	require.Len(t, capturedT.captured.params, 4)
	require.Equal(t, "resp_value1", getParameterValue(capturedT.captured.params, "resp_key1"))
	require.Equal(t, "resp_value2", getParameterValue(capturedT.captured.params, "resp_key2"))
	require.Equal(t, "resp_value3", getParameterValue(capturedT.captured.params, "resp_key3"))
	require.Equal(t, "resp_value4", getParameterValue(capturedT.captured.params, "resp_key4"))
}

func TestResponseInformationT(t *testing.T) {
	t.Parallel()
	var capturedResponse *http.Response

	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://httpbin.org/get"),
			},
		},
		ResponseInformationT: []ResponseInformationT{
			func(t T, resp *http.Response) ([]*allure.Parameter, error) {
				capturedResponse = resp
				return allure.NewParameters("t_response_callback", "response_works"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	require.NotNil(t, capturedResponse)
	require.Equal(t, http.StatusOK, capturedResponse.StatusCode)
	// Verify that the parameter from ResponseInformationT was captured
	require.NotEmpty(t, capturedT.captured.params)
	require.Len(t, capturedT.captured.params, 1)
	require.Equal(t, "response_works", getParameterValue(capturedT.captured.params, "t_response_callback"))
}

func TestResponseInformationT_Error(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/test"),
			},
		},
		ResponseInformationT: []ResponseInformationT{
			func(t T, resp *http.Response) ([]*allure.Parameter, error) {
				return nil, errors.New("t_response_callback error")
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	// Errors in ResponseInformationT are logged but do not fail the test
	// The request is still executed successfully
	require.Empty(t, results.GetErrors())
	// No parameters should be captured when the callback returns an error
	require.Empty(t, capturedT.captured.params)
}

func TestResponseInformationT_MultipleCallbacks(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/test"),
			},
		},
		ResponseInformationT: []ResponseInformationT{
			func(t T, resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("t_resp_key1", "t_resp_value1"), nil
			},
			func(t T, resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("t_resp_key2", "t_resp_value2"), nil
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
	require.Equal(t, "t_resp_value1", getParameterValue(capturedT.captured.params, "t_resp_key1"))
	require.Equal(t, "t_resp_value2", getParameterValue(capturedT.captured.params, "t_resp_key2"))
}

func TestResponseInformationAndResponseInformationT_Together(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/test"),
			},
		},
		ResponseInformation: []ResponseInformation{
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("resp_key", "resp_value"), nil
			},
		},
		ResponseInformationT: []ResponseInformationT{
			func(t T, resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("resp_t_key", "resp_t_value"), nil
			},
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Should capture parameters from both ResponseInformation and ResponseInformationT
	require.Len(t, capturedT.captured.params, 2)
	require.Equal(t, "resp_value", getParameterValue(capturedT.captured.params, "resp_key"))
	require.Equal(t, "resp_t_value", getParameterValue(capturedT.captured.params, "resp_t_key"))
}

// Tests for ResponseBaseInformation
func TestResponseBaseInformation(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 201,
		Header:     make(http.Header),
		Body:       nil,
	}

	params, err := ResponseBaseInformation(resp)

	require.NoError(t, err)
	require.NotNil(t, params)
	require.Len(t, params, 1)
	require.Equal(t, "201", getParameterValue(params, "response_code"))
}

// Tests for ResponseHeadersInformation
func TestResponseHeadersInformation(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": {"application/json"},
			"X-Custom":     {"custom-value"},
		},
		Body: nil,
	}

	params, err := ResponseHeadersInformation(resp)

	require.NoError(t, err)
	require.NotNil(t, params)
	require.Len(t, params, 1)

	headersValue := getParameterValue(params, "response_headers")
	require.NotNil(t, headersValue)
	headersStr := headersValue.(string)
	require.Contains(t, headersStr, "Content-Type")
	require.Contains(t, headersStr, "X-Custom")
}

func TestResponseHeadersInformation_NoHeaders(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       nil,
	}

	params, err := ResponseHeadersInformation(resp)

	require.NoError(t, err)
	require.NotNil(t, params)

	headersValue := getParameterValue(params, "response_headers")
	require.NotNil(t, headersValue)
	// Empty headers should still produce a valid JSON object
	headersStr := headersValue.(string)
	require.Equal(t, "{}", headersStr)
}

// Tests for ResponseBodyInformationT
func TestResponseBodyInformationT_WithJSONBody(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": {"application/json"},
		},
		Body: &readCloser{data: []byte(`{"key":"value","number":123}`)},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	// ResponseBodyInformationT adds attachments, not parameters
	// We need to verify attachments were added
	_, err := ResponseBodyInformationT(capturedT, resp)

	require.NoError(t, err)
	// Check that attachment was added
	require.Len(t, capturedT.captured.attachments, 1)
	attachment := getAttachment(capturedT.captured.attachments, "response")
	require.NotNil(t, attachment)
	require.Equal(t, allure.JSON, attachment.Type)
}

func TestResponseBodyInformationT_WithTextBody(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": {"text/plain"},
		},
		Body: &readCloser{data: []byte("plain text response")},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	_, err := ResponseBodyInformationT(capturedT, resp)

	require.NoError(t, err)
	require.Len(t, capturedT.captured.attachments, 1)
	attachment := getAttachment(capturedT.captured.attachments, "response")
	require.NotNil(t, attachment)
	require.Equal(t, allure.Text, attachment.Type)
}

func TestResponseBodyInformationT_EmptyBody(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 204,
		Header:     make(http.Header),
		Body:       nil,
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	params, err := ResponseBodyInformationT(capturedT, resp)

	require.NoError(t, err)
	// Should return nil when body is empty
	require.Nil(t, params)
	// No attachments should be added for empty body
	require.Empty(t, capturedT.captured.attachments)
}

func TestResponseBodyInformationT_NilBody(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       nil,
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	params, err := ResponseBodyInformationT(capturedT, resp)

	require.NoError(t, err)
	// Should return nil when body is nil
	require.Nil(t, params)
	// No attachments should be added for nil body
	require.Empty(t, capturedT.captured.attachments)
}

func TestResponseBodyInformationT_CustomMimeType(t *testing.T) {
	t.Parallel()
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": {"text/html; charset=utf-8"},
		},
		Body: &readCloser{data: []byte("<html><body>test</body></html>")},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)

	_, err := ResponseBodyInformationT(capturedT, resp)

	require.NoError(t, err)
	require.Len(t, capturedT.captured.attachments, 1)
	attachment := getAttachment(capturedT.captured.attachments, "response")
	require.NotNil(t, attachment)
	require.Equal(t, allure.MimeType("text/html; charset=utf-8"), attachment.Type)
}

func TestResponseInformation_IntegratesWithResponseBodyInformationT(t *testing.T) {
	t.Parallel()
	test := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(testServerAddress + "/with_body"),
			},
		},
		ResponseInformation: []ResponseInformation{
			func(resp *http.Response) ([]*allure.Parameter, error) {
				return allure.NewParameters("response_code", resp.StatusCode), nil
			},
		},
		ResponseInformationT: []ResponseInformationT{
			ResponseBodyInformationT,
		},
	}

	allureT := createAllureT(t)
	capturedT := newCaptureT(allureT)
	results := test.Execute(context.Background(), capturedT)

	require.NotNil(t, results)
	require.Empty(t, results.GetErrors())
	// Verify response code parameter was captured
	require.NotEmpty(t, capturedT.captured.params)
	require.Equal(t, "200", getParameterValue(capturedT.captured.params, "response_code"))
	// Verify attachment was added
	require.Len(t, capturedT.captured.attachments, 1)
	attachment := getAttachment(capturedT.captured.attachments, "response")
	require.NotNil(t, attachment)
}

// Helper types for testing

type readCloser struct {
	data []byte
	pos  int
}

func (r *readCloser) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	copy(p, r.data[r.pos:])
	n = len(r.data) - r.pos
	r.pos = len(r.data)
	return n, nil
}

func (r *readCloser) Close() error {
	return nil
}
