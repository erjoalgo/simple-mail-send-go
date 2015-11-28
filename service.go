package main

import (
	"fmt"
	"io"
	"net/http"
)

type MailService struct {
	EndpointUrl, ServiceName string
	// WasSuccess(http.Response) error

	//convert generic SendMailRequest to (data, content-type)
	TranslateMailRequest func(SendMailRequest, ApiKey) (io.Reader, string, error)

	//service-specific http.Request modifications
	PostRequestHook func(*http.Request, SendMailRequest, ApiKey) error
}

// Any type of credential-related information needed by a service
type ApiKey interface{}

func (s MailService) Send(mr SendMailRequest) (*http.Response, error) {
	if apiKey := mr.ApiKeyFor(s); apiKey == nil {
		return nil, fmt.Errorf(
			"no apiKey provided in mail request: %v for service %s\n",
			mr, s)
	} else if data, contentType, err := s.TranslateMailRequest(mr, apiKey); err != nil {
		return nil, fmt.Errorf(
			"error translating mail request: %v",
			mr, err)
	} else if req, err := http.NewRequest("POST", s.EndpointUrl, data); err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	} else {
		req.Header.Set("Content-Type", contentType)
		if s.PostRequestHook != nil {
			if err := s.PostRequestHook(req, mr, apiKey); err != nil {
				return nil, fmt.Errorf("error in post-request hook: %v", err)
			}
		}
		if resp, err := http.DefaultClient.Do(req); err != nil {
			return nil, fmt.Errorf("error executing request: %v", err)
		} else {
			return resp, nil
		}
	}
}

// Determine whether Send request was successful in a possibly service-specific way
/*
func (s MailService) WasSuccess(r http.Response) error {
	// return nil
	if r.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("non 200 status code: %v", r)
	}
}
*/

const MAX_RETRY = 10

// Try sending using several services, several times
func SendRetry(mr SendMailRequest, maxRetry int) (resp *http.Response, errors []error, success bool) {
	for retriesLeft := maxRetry; retriesLeft > 0; retriesLeft-- {
		for _, service := range MailServices {
			if resp, err := service.Send(mr); err != nil {
				errors = append(errors, err)
			} else {
				return resp, errors, true
			}
		}
	}
	return
}
