package cute

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/cute/internal/utils"
	"io"
	"net/http"
	"strings"
)

type ResponseInformation func(*http.Response) ([]*allure.Parameter, error)
type ResponseInformationT func(T, *http.Response) ([]*allure.Parameter, error)

func ResponseBaseInformation(resp *http.Response) ([]*allure.Parameter, error) {
	return allure.NewParameters(
		"response_code", resp.StatusCode,
	), nil
}

func ResponseHeadersInformation(resp *http.Response) ([]*allure.Parameter, error) {
	// Do not change to JSONMarshaler
	// In this case we can keep default for keep JSON, independence from JSONMarshaler
	headers, err := utils.ToJSON(resp.Header)
	if err != nil {
		return nil, err
	}
	return allure.NewParameters("response_headers", headers), nil
}

func ResponseBodyInformationT(t T, resp *http.Response) ([]*allure.Parameter, error) {
	if resp.Body == nil {
		return nil, nil
	}

	var (
		saveBody io.ReadCloser
		err      error
	)

	saveBody, resp.Body, err = utils.DrainBody(resp.Body)
	// if could not get body from response, no add to allure
	if err != nil {
		return nil, err
	}

	body, err := utils.GetBody(saveBody)
	// if could not get body from response, no add to allure
	if err != nil {
		return nil, err
	}

	// if body is empty - skip
	if len(body) == 0 {
		return nil, nil
	}

	responseType := allure.Text

	if _, ok := resp.Header["Content-Type"]; ok {
		if len(resp.Header["Content-Type"]) > 0 {
			if strings.Contains(resp.Header["Content-Type"][0], "application/json") {
				responseType = allure.JSON
			} else {
				responseType = allure.MimeType(resp.Header["Content-Type"][0])
			}
		}
	}

	if responseType == allure.JSON {
		body, _ = utils.PrettyJSON(body)
	}

	t.WithAttachments(allure.NewAttachment("response", responseType, body))

	return nil, nil
}
