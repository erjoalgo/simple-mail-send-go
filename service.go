package main

import (
	"fmt"
	"io"
	"net/http"
)

type MailService struct {
	ServiceName string

	EndpointUrl func(SendMailRequest, Credentials) (string, error)

	//convert generic SendMailRequest to (data, content-type)
	TranslateMailRequest func(SendMailRequest, Credentials) (io.Reader, string, error)

	//service-specific http.Request modifications
	PostRequestHook func(*http.Request, SendMailRequest, Credentials) error

	WasSuccessFunc      func(*http.Response) error
	RequiredCredentials []string

	DocUrl string
}

func (s MailService) Send(mr SendMailRequest) (*http.Response, error) {
	if credentials := mr.CredentialsFor(s); credentials == nil {
		return nil, fmt.Errorf(
			"no credentials provided in mail request: %v for service %s\n",
			mr, s)
	} else if data, contentType, err := s.TranslateMailRequest(mr, credentials); err != nil {
		return nil, fmt.Errorf(
			"error translating mail request: %v for: %v", err, mr)
	} else if endpointUrl, err := s.EndpointUrl(mr, credentials); err != nil {
		return nil, fmt.Errorf("error obtaining endpoint url: %v", err)
	} else if req, err := http.NewRequest("POST", endpointUrl, data); err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	} else {
		req.Header.Set("Content-Type", contentType)
		if s.PostRequestHook != nil {
			if err := s.PostRequestHook(req, mr, credentials); err != nil {
				return nil, fmt.Errorf("error in post-request hook: %v", err)
			}
		}
		if resp, err := http.DefaultClient.Do(req); err != nil {
			return nil, fmt.Errorf("error executing request: %v", err)
		} else {
			return resp, s.WasSuccess(resp)
		}
	}
}

// Determine whether Send request was successful in a possibly service-specific way
func (s MailService) WasSuccess(r *http.Response) error {
	// return nil
	if s.WasSuccessFunc != nil {
		return s.WasSuccessFunc(r)
	} else if r.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("non 200 status code: %v", r)
	}
}

// Number of retries for each service
const MAX_RETRY = 3

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
