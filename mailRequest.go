package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type SendMailRequest struct {
	FromEmail   string
	ToEmails    []string
	Subject     string
	MessageText string
	// some services require several credentials
	Credentials    map[string]Credentials `json:"-"`
	RawCredentials map[string]interface{} `json:"Credentials"`
}

func (mr SendMailRequest) CredentialsFor(s MailService) Credentials {
	return mr.Credentials[strings.ToLower(s.ServiceName)]
}

// Any type of credential-related information needed by a service
type Credentials map[string]interface{}

func (c Credentials) Get(key string) (string, bool) {
	if val, contains := c[key]; !contains {
		return "", false
	} else if valString, ok := val.(string); !ok {
		return "", false
	} else {
		return valString, true
	}
}

func (mr SendMailRequest) Validate() error {
	if mr.FromEmail == "" {
		return fmt.Errorf("FromEmail is required")
	} else if len(mr.ToEmails) == 0 {
		return fmt.Errorf("ToEmails is required")
	} else if len(mr.Credentials) == 0 {
		return fmt.Errorf("Credentials required")
	} else {
		return nil
	}
}
func NewMailRequestJson(jsonRequest io.Reader) (mr SendMailRequest, err error) {
	if err = json.NewDecoder(jsonRequest).Decode(&mr); err != nil {
		return
	} else {
		mr.Credentials = make(map[string]Credentials)
		for key, value := range mr.RawCredentials {
			if valueMap, ok := value.(map[string]interface{}); !ok {
				return mr, fmt.Errorf(
					"all credentials must be key-value pairs: %v",
					value)
			} else {
				mr.Credentials[key] = Credentials(valueMap)
			}
		}
	}
	return
}
