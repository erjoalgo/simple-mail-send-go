package main

import (
	"fmt"
	"net/url"
	"strings"
)

type SendMailRequest struct {
	FromEmail   string
	ToEmails    []string
	Subject     string
	MessageText string
	// some services require several credentials
	ApiKeys map[string]ApiKey
}

func (mr SendMailRequest) ApiKeyFor(s MailService) ApiKey {
	return mr.ApiKeys[strings.ToLower(s.ServiceName)]
}

var REQUIRED_FIELDS = []string{"FromEmail", "ToEmails", "ApiKeys"}

func ParseMailRequest(values url.Values) (SendMailRequest, error) {
	for _, field := range REQUIRED_FIELDS {
		if list, contains := values[field]; !contains && len(list) == 0 {
			return SendMailRequest{}, fmt.Errorf("required field %s missing or empty in request", field)
		}
	}
	return SendMailRequest{
		FromEmail:   values["FromEmail"][0],
		ToEmails:    values["ToEmails"],
		Subject:     values["Subject"][0],
		MessageText: values["MessageText"][0],
		// TODO JSON
		// ApiKeys:     values["ApiKeys"],
	}, nil
}
