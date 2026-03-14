package cute

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/cute/internal/utils"
	"io"
	"moul.io/http2curl/v2"
	"net/http"
)

type RequestInformation func(*http.Request) ([]*allure.Parameter, error)
type RequestInformationT func(T, *http.Request) ([]*allure.Parameter, error)

func RequestCurlInformation(req *http.Request) ([]*allure.Parameter, error) {
	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	return allure.NewParameters("curl", curl.String()), nil
}

func RequestHTTPBaseInformation(req *http.Request) ([]*allure.Parameter, error) {
	return allure.NewParameters(
		"method", req.Method,
		"host", req.Host,
	), nil
}

func RequestHTTPHeadersInformation(req *http.Request) ([]*allure.Parameter, error) {
	// Do not change to JSONMarshaler
	// In this case we can keep default for keep JSON, independence from JSONMarshaler
	headers, err := utils.ToJSON(req.Header)
	if err != nil {
		return nil, err
	}
	return allure.NewParameters("headers", headers), nil
}

func RequestHTTPBodyInformation(req *http.Request) ([]*allure.Parameter, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return nil, nil
	}
	var (
		saveBody io.ReadCloser
		err      error
	)
	saveBody, req.Body, err = utils.DrainBody(req.Body)
	if err != nil {
		return nil, err
	}

	body, err := utils.GetBody(saveBody)
	if err != nil {
		return nil, err
	}

	if len(body) != 0 {
		return nil, nil
	}
	return allure.NewParameters("body", string(body)), nil
}
