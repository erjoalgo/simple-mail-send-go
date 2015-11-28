package main

// This file provides various service implementations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_URLENCODED = "application/x-www-form-urlencoded"

var SendGrid = MailService{
	ServiceName: "SendGrid",
	EndpointUrl: "https://api.sendgrid.com/api/mail.send.json?",
	TranslateMailRequest: func(mr SendMailRequest, apiKey ApiKey) (io.Reader, string, error) {
		if key, ok := apiKey.(map[string]string); !ok {
			return nil, "", fmt.Errorf("sendgrid key must be a map, not %v", key)
		} else if key["key"] == "" || key["user"] == "" {
			return nil, "", fmt.Errorf("sendgrid key must contain user, key %v", key)
		} else {
			data := url.Values{
				"api_key":  {key["key"]},
				"api_user": {key["user"]},
				"from":     {mr.FromEmail},
				"to":       {mr.ToEmails[0]}, //TODO
				"subject":  {mr.Subject},
				"text":     {mr.MessageText},
				"headers":  {"null"},
				"html":     {""},
			}
			fmt.Printf("%s\n", data.Encode())
			return strings.NewReader(data.Encode()), CONTENT_TYPE_URLENCODED, nil
		}
	},
}

var MailGun = MailService{
	ServiceName: "MailGun",
	EndpointUrl: "https://api.mailgun.net/v3/sandbox4a89fc7b8f8f4a80a22bf7835a3ff6cb.mailgun.org/messages?",
	TranslateMailRequest: func(mr SendMailRequest, apiKey ApiKey) (io.Reader, string, error) {
		data := url.Values{
			// "from":    {"Excited User <mailgun@sandbox4a89fc7b8f8f4a80a22bf7835a3ff6cb.mailgun.org>"},
			"from": {fmt.Sprintf("<%s>", mr.FromEmail)},
			// TODO send multiple emails
			"to":      {mr.ToEmails[0]},
			"subject": {mr.Subject},
			"text":    {mr.MessageText},

			"headers": {"null"},
			"html":    {"null"},
		}

		return strings.NewReader(data.Encode()), CONTENT_TYPE_URLENCODED, nil
	},
	PostRequestHook: func(req *http.Request, mr SendMailRequest, key ApiKey) error {
		if key, ok := key.(string); !ok {
			return fmt.Errorf("mailgun requires an api key", key)
		} else {
			req.SetBasicAuth("api", key)
			return nil
		}
	},
}

var Mandrill = MailService{
	ServiceName: "Mandrill",
	EndpointUrl: "https://mandrillapp.com/api/1.0/messages/send.json",
	TranslateMailRequest: func(mr SendMailRequest, apiKey ApiKey) (io.Reader, string, error) {
		if key, ok := apiKey.(string); !ok {
			return nil, "", fmt.Errorf("mandrill requires a key")
		} else {
			data := map[string]interface{}{
				"key": key,
				"message": map[string]interface{}{
					"from_email": mr.FromEmail,
					"text":       mr.MessageText,
					"subject":    mr.Subject,
					// "from_name":  Example Name",
					"to": []interface{}{map[string]interface{}{
						// "email": "erjoalgo@gmail.com",
						"email": mr.ToEmails[0],
						// "name":  "Recipient Name",
						"type": "to",
					}},
				},
			}
			if data, err := json.Marshal(data); err != nil {
				return nil, "", err
			} else {
				return bytes.NewReader(data), CONTENT_TYPE_JSON, nil
			}
		}
	},
}

var Amazon = MailService{
	ServiceName: "Amazon",
	EndpointUrl: "https://email.us-east-1.amazonaws.com",
	TranslateMailRequest: func(mr SendMailRequest, apiKey ApiKey) (io.Reader, string, error) {
		if key, ok := apiKey.(string); !ok {
			return nil, "", fmt.Errorf("amazon requires a key")
		} else {
			data := url.Values{
				"AWSAccessKeyId": {key},
				"Action":         {"SendEmail"},
				"Destination.ToAddress.member.1": {"erjoalgo@gmail.com"},
				"Source":                 {"erjoalgo@gmail.com"},
				"Message.Subject.Data":   {mr.Subject},
				"Message.Body.Text.Data": {mr.MessageText},
				"Version":                {"2010-12-01"},
			}
			return strings.NewReader(data.Encode()), CONTENT_TYPE_URLENCODED, nil
		}
	},
	PostRequestHook: func(req *http.Request, mr SendMailRequest, key ApiKey) error {
		req.Header.Set("X-Amzn-Authorization", "AWS3-HTTPS")
		return nil
	},
}

var MailServices = []MailService{SendGrid, MailGun, Mandrill, Amazon}
