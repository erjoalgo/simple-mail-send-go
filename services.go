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
	EndpointUrl: func(mr SendMailRequest, credentials Credentials) (string, error) {
		return "https://api.sendgrid.com/api/mail.send.json?", nil
	},
	TranslateMailRequest: func(mr SendMailRequest, credentials Credentials) (io.Reader, string, error) {
		key, okKey := credentials.Get("passwd")
		user, okUser := credentials.Get("user")
		if !okUser || !okKey {
			return nil, "", fmt.Errorf("sendgrid credentials must contain 'user', 'passwd' %v", credentials)
		} else {
			data := url.Values{
				"api_key":  {key},
				"api_user": {user},
				"from":     {mr.FromEmail},
				"to":       {mr.ToEmails[0]}, //TODO
				"subject":  {mr.Subject},
				"text":     {mr.MessageText},
				"headers":  {"null"},
				"html":     {"null"}, //unfortunately required
			}
			return strings.NewReader(data.Encode()), CONTENT_TYPE_URLENCODED, nil
		}
	},
	RequiredCredentials: []string{"passwd", "user"},
	DocUrl:              "https://sendgrid.com/docs/API_Reference/Web_API/mail.html",
}

var MailGun = MailService{
	ServiceName: "MailGun",
	TranslateMailRequest: func(mr SendMailRequest, credentials Credentials) (io.Reader, string, error) {
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
	PostRequestHook: func(req *http.Request, mr SendMailRequest, key Credentials) error {
		if key, ok := key.Get("apikey"); !ok {
			return fmt.Errorf("mailgun requires an api key", key)
		} else {
			req.SetBasicAuth("api", key)
			return nil
		}
	},
	EndpointUrl: func(mr SendMailRequest, credentials Credentials) (string, error) {
		if domain, ok := credentials.Get("domain"); !ok {
			return "", fmt.Errorf("mailgun requires domain in api key: %v", credentials)
		} else {
			return "https://api.mailgun.net/v3/" + domain + "/messages?", nil
		}
	},
	RequiredCredentials: []string{"domain", "apikey"},
	DocUrl:              "https://documentation.mailgun.com/quickstart.html#sending-messages",
}

var Mandrill = MailService{
	ServiceName: "Mandrill",
	EndpointUrl: func(mr SendMailRequest, credentials Credentials) (string, error) {
		return "https://mandrillapp.com/api/1.0/messages/send.json", nil
	},
	TranslateMailRequest: func(mr SendMailRequest, credentials Credentials) (io.Reader, string, error) {
		if apikey, ok := credentials.Get("apikey"); !ok {
			return nil, "", fmt.Errorf("mandrill requires an apikey")
		} else {
			data := map[string]interface{}{
				"key": apikey,
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
	RequiredCredentials: []string{"apikey"},
	DocUrl:              "https://mandrillapp.com/api/docs/messages.JSON.html#method-send",
}

var Amazon = MailService{
	ServiceName: "Amazon",
	EndpointUrl: func(mr SendMailRequest, credentials Credentials) (string, error) {
		return "https://email.us-east-1.amazonaws.com", nil
	},
	TranslateMailRequest: func(mr SendMailRequest, credentials Credentials) (io.Reader, string, error) {
		if key, ok := credentials.Get("access-key-id"); !ok {
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
	PostRequestHook: func(req *http.Request, mr SendMailRequest, key Credentials) error {
		req.Header.Set("X-Amzn-Authorization", "AWS3-HTTPS")
		return nil
	},
	RequiredCredentials: []string{"access-key-id", "secret-access-key"},
	DocUrl:              "http://docs.aws.amazon.com/ses/latest/APIReference/API_SendEmail.html",
}

var MailServices = []MailService{SendGrid, MailGun, Mandrill, Amazon}
