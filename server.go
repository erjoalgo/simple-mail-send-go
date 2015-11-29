package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ServeTLS(address string, certPem string, keyPem string) error {
	http.HandleFunc("/ok", okHandler)
	http.HandleFunc("/send", mailRequestHandler)
	http.HandleFunc("/credentials", credentialListerHandler)
	http.HandleFunc("/sampleRequest", exampleRequestHandler)
	return http.ListenAndServeTLS(address, certPem, keyPem, nil)
}

func okHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req, "")
	w.WriteHeader(200)
	fmt.Fprintf(w, "OK")
}

func exampleRequestHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, sampleRequest)
}

func credentialListerHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	logRequest(req, "")
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String()
	logRequest(req, s)
	if mr, err := NewMailRequestJson(strings.NewReader(s)); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error parsing mail request: %s", err)
	} else if resp, errs, succ := SendRetry(mr, MAX_RETRY); !succ {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to deliver mail after max retries: %s", errs)
	} else {
		w.WriteHeader(200)
		fmt.Fprintf(w, "success: %v", resp)
	}
}

func logRequest(req *http.Request, body string) {
	log.Printf("NEW REQUEST:")
	log.Printf("Method: %s", req.Method)
	log.Printf("URL: %s", req.URL)
	log.Printf("Proto: %s", req.Proto)
	log.Printf("ProtoMajor: %d", req.ProtoMajor)
	log.Printf("ProtoMinor: %d", req.ProtoMinor)
	log.Printf("Header: %#v", req.Header)
	log.Printf("Body: %s", body)
	log.Printf("ContentLength: %d", req.ContentLength)
	log.Printf("TransferEncoding: %s", req.TransferEncoding)
	log.Printf("Host: %s", req.Host)
	log.Printf("Form: %s", req.Form)
	log.Printf("PostForm: %s", req.PostForm)
	log.Printf("MultipartForm: %s", req.MultipartForm)
	log.Printf("Trailer: %s", req.Trailer)
	log.Printf("RemoteAddr: %s", req.RemoteAddr)
	log.Printf("RequestURI: %s", req.RequestURI)
	log.Printf("TLS: %#v", req.TLS)
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
