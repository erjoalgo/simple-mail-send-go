package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func init() {
	ServeTLS("localhost:6700", "cert.pem", "key.pem")
}

var mailRequest SendMailRequest
var jsonMR []byte

func TestJSON(t *testing.T) {
	filename := "mailRequest.json"
	if json, err := ioutil.ReadFile(filename); err != nil {
		t.Errorf("error reading credentials file: %s, %s", filename, err)
	} else if mr, err := NewMailRequestJson(bytes.NewReader(json)); err != nil {
		t.Errorf("error parsing mail request: %s ", err)
	} else {
		t.Logf("mr: %#v", mr)
		if err := mr.Validate(); err != nil {
			t.Errorf("invalid mail request: %s", err)
		} else {
			mailRequest = mr
			jsonMR = json
		}
	}
}

func TestPost(t *testing.T) {
	http.Post("localhost:6700", CONTENT_TYPE_JSON, bytes.NewReader(jsonMR))
}
func TestSend(t *testing.T) {
	var mr = mailRequest
	origSubject := mr.Subject
	t.Logf("Using mail request: %#v", mr)
	for _, service := range MailServices {
		mr.Subject = origSubject + service.ServiceName
		t.Logf("using service: %s \n", service.ServiceName)
		if resp, err := service.Send(mr); err != nil {
			t.Errorf("err: %#v\n", err)
		} else {
			t.Errorf("%#v\n", resp)
		}
	}
}

// Local Variables:
// compile-cmd: "go test"
// End:
