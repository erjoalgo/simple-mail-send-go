package main

import (
	"fmt"
	"log"
	"net/http"
)

func ServeTLS(address string, certPem string, keyPem string) error {
	http.HandleFunc("/ok", okHandler)
	http.HandleFunc("/send", mailRequestHandler)
	http.HandleFunc("/credentials", credentialListerHandler)
	http.HandleFunc("/sampleRequest", exampleRequestHandler)
	return http.ListenAndServeTLS(address, certPem, keyPem, nil)
}

func okHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "OK")
}

func exampleRequestHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, sampleRequest)
}

func credentialListerHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	for _, service := range MailServices {
		fmt.Fprintf(w, "%s\n\trequires: %s\n\tmore info: %s\n",
			service.ServiceName, service.RequiredCredentials, service.DocUrl)
	}
}
func mailRequestHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("$+v", r)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
		}
	}()
	if mr, err := NewMailRequestJson(req.Body); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error parsing mail request: %s", err)
	} else if resp, errs, succ, service := SendRetry(mr, MAX_RETRY); !succ {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to deliver mail after max retries: %s", errs)
	} else {
		w.WriteHeader(200)
		fmt.Fprintf(w, "success with %s: %v", service.ServiceName, resp)
	}
}

const sampleRequest = `{
    "FromEmail": "me@example.com",
    "ToEmails": ["you@example.com"],
    "Subject":  "hola from simple-mail-send-go",
    "MessageText": "message text",
    "Credentials": {
	"mandrill": {"apikey":"<key>"},
	"mailgun": {"apikey":"<key>",
		    "domain":"<domain>"},
	"sendgrid": {
	    "user": "<user>", 
	    "passwd":  "<pass>"
	},
	"amazon": {"access-key-id":"<key>", "secret-access-key":"<key>"}
    }
}`
