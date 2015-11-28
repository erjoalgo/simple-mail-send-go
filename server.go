package main

import (
	"fmt"
	"net/http"
)

func Serve(address string) error {
	r := http.NewServeMux()
	r.HandleFunc("/ok", okHandler)
	r.HandleFunc("/", mailRequestHandler)
	return http.ListenAndServe(address, r)
}

func okHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "OK")
}

func mailRequestHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// log.Critical("$+v", r)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
		}
	}()
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error parsing post form: %v", err)
	} else if mr, err := ParseMailRequest(req.PostForm); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error parsing send mail request: %v", err)
	} else if resp, errs, succ := SendRetry(mr, MAX_RETRY); !succ {
		w.WriteHeader(400)
		fmt.Fprintf(w, "unable to deliver mail after max retries: %s", errs)
	} else {
		w.WriteHeader(200)
		fmt.Fprintf(w, "success: %s", resp)
	}
}
