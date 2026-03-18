package cute

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	cuteErrors "github.com/slowaner/cute/errors"
	"github.com/slowaner/cute/internal/utils"
)

func (it *Test) makeRequest(t internalT, req *http.Request) (*http.Response, []error) {
	var (
		delay       = defaultDelayRepeat
		countRepeat = 1

		resp  *http.Response
		err   error
		scope = make([]error, 0)
	)

	if it.Request.Retry.Delay != 0 {
		delay = it.Request.Retry.Delay
	}

	if it.Request.Retry.Count != 0 {
		countRepeat = it.Request.Retry.Count
	}

	for i := 1; i <= countRepeat; i++ {
		it.executeWithStep(t, it.createTitle(i, countRepeat, req), func(t T) []error {
			resp, err = it.doRequest(t, req)
			if err != nil {
				if it.Request.Retry.Broken {
					err = wrapBrokenError(err)
				}

				if it.Request.Retry.Optional {
					err = wrapOptionalError(err)
				}

				return []error{err}
			}

			return nil
		})

		if err == nil {
			break
		}

		scope = append(scope, err)

		if i != countRepeat {
			time.Sleep(delay)
		}
	}

	return resp, scope
}

func (it *Test) doRequest(t T, baseReq *http.Request) (*http.Response, error) {
	// copy request, because body can be read once
	req, err := copyRequest(baseReq.Context(), baseReq)
	if err != nil {
		return nil, cuteErrors.NewCuteError("[Internal] Could not copy request", err)
	}

	resp, httpErr := it.httpClient.Do(req)

	// if the timeout is triggered, we properly log the timeout error on allure and in traces
	if errors.Is(httpErr, context.DeadlineExceeded) {
		// Add information (method, host, curl) about request to Allure step
		// should be after httpClient.Do and from resp.Request, because in roundTripper request may be changed
		if addErr := it.addInformationRequest(t, req); addErr != nil {
			// Ignore err return, because it's connected with test logic
			it.Error(t, "Could not log information about request. error %v", addErr)
		}

		return nil, cuteErrors.NewEmptyAssertError(
			"Request timeout",
			fmt.Sprintf("expected request to be completed in %v, but was not", it.Expect.ExecuteTime))
	}

	// http client has case when it returns response and error in one time
	// we have to check this case
	if resp == nil {
		if httpErr != nil {
			return nil, cuteErrors.NewCuteError("[HTTP] Could not do request", httpErr)
		}

		// if response is nil, we can't get information about request and response
		return nil, cuteErrors.NewCuteError("[HTTP] Response is nil", httpErr)
	}

	// BAD CODE. Need to copy body, because we can't read body again from resp.Request.Body. Problem is io.Reader
	resp.Request.Body, baseReq.Body, err = utils.DrainBody(baseReq.Body)
	if err != nil {
		it.Error(t, "Could not drain body from baseReq.Body. error %v", err)
		// Ignore err return, because it's connected with test logic
	}

	// Add information (method, host, curl) about request to Allure step
	// should be after httpClient.Do and from resp.Request, because in roundTripper request may be changed
	if addErr := it.addInformationRequest(t, resp.Request); addErr != nil {
		it.Error(t, "Could not log information about request. error %v", addErr)
		// Ignore err return, because it's connected with test logic
	}

	if httpErr != nil {
		return nil, cuteErrors.NewCuteError("[HTTP] Could not do request", httpErr)
	}

	// Add information (code, body, headers) about response to Allure step
	if addErr := it.addInformationResponse(t, resp); addErr != nil {
		// Ignore err return, because it's connected with test logic
		it.Error(t, "Could not log information about response. error %v", addErr)
	}

	if validErr := it.validateResponseCode(resp); validErr != nil {
		return resp, validErr
	}

	return resp, nil
}

func (it *Test) validateResponseCode(resp *http.Response) error {
	if it.Expect.Code != 0 && it.Expect.Code != resp.StatusCode {
		return cuteErrors.NewAssertError(
			"Assert response code",
			fmt.Sprintf("Response code expect %v, but was %v", it.Expect.Code, resp.StatusCode),
			resp.StatusCode,
			it.Expect.Code)
	}

	return nil
}

func (it *Test) addInformationRequest(t T, req *http.Request) error {
	// sanitize in any way
	//  FIXME: why? this method must only pass parameters into T.
	if it.RequestSanitizer != nil {
		it.RequestSanitizer(req)
	}

	if len(it.RequestInformation)+len(it.RequestInformationT) == 0 {
		return nil
	}

	allParameters := make([]*allure.Parameter, 0, len(it.RequestInformation)+len(it.RequestInformationT))
	for _, information := range it.RequestInformation {
		parameters, err := information(req)
		if err != nil {
			return err
		}
		allParameters = append(allParameters, parameters...)
	}
	for _, information := range it.RequestInformationT {
		parameters, err := information(t, req)
		if err != nil {
			return err
		}
		allParameters = append(allParameters, parameters...)
	}

	if len(allParameters) == 0 {
		return nil
	}

	t.WithParameters(allParameters...)

	return nil
}

func copyRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	var (
		err error

		clone = req.Clone(ctx)
	)

	req.Body, clone.Body, err = utils.DrainBody(req.Body)
	if err != nil {
		return nil, err
	}

	return clone, nil
}

func (it *Test) addInformationResponse(t T, response *http.Response) error {
	// sanitize in any way
	//  FIXME: why? this method must only pass parameters into T.
	if it.ResponseSanitizer != nil {
		it.ResponseSanitizer(response)
	}

	if len(it.ResponseInformation)+len(it.ResponseInformationT) == 0 {
		return nil
	}

	allParameters := make([]*allure.Parameter, 0, len(it.ResponseInformation)+len(it.ResponseInformationT))
	for _, information := range it.ResponseInformation {
		parameters, err := information(response)
		if err != nil {
			return err
		}
		allParameters = append(allParameters, parameters...)
	}
	for _, information := range it.ResponseInformationT {
		parameters, err := information(t, response)
		if err != nil {
			return err
		}
		allParameters = append(allParameters, parameters...)
	}

	if len(allParameters) == 0 {
		return nil
	}

	t.WithParameters(allParameters...)

	return nil
}

func (it *Test) createTitle(try, countRepeat int, req *http.Request) string {
	toProcess := req

	// We have to execute sanitizer hook because
	// we need to log it and it can contain sensitive data
	if it.RequestSanitizer != nil {
		clone, err := copyRequest(req.Context(), req)

		// ignore error, because we want to log request
		// and it does not matter if we can copy request
		if err == nil {
			it.RequestSanitizer(clone)

			toProcess = clone
		}
	}

	title := toProcess.Method + " " + toProcess.URL.String()

	if countRepeat == 1 {
		return title
	}

	return fmt.Sprintf("[%v/%v] %v", try, countRepeat, title)
}
