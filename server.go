package main

import (
	"fmt"
	"net/http"
)

func ServeTLS(address string, certPem string, keyPem string) error {
	http.HandleFunc("/ok", okHandler)
	http.HandleFunc("/send", mailRequestHandler)
	http.HandleFunc("/credentials", credentialListerHandler)
	return http.ListenAndServeTLS(address, certPem, keyPem, nil)
}

func okHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "OK")
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
			// log.Critical("$+v", r)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
		}
	}()
	if mr, err := NewMailRequestJson(req.Body); err != nil {
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
