package main

import (
	"flag"
	"log"
)

func main() {
	var certPem string
	var keyPem string
	var address string
	flag.StringVar(&certPem, "certPem", "", "cert.pem file location")
	flag.StringVar(&keyPem, "keyPem", "", "key.pem file location")
	flag.StringVar(&address, "address", "", "listen address")
	flag.Parse()
	if address == "" {
		log.Fatal("address argument required")
	} else if keyPem == "" || certPem == "" {
		log.Fatal("cert, key files required")
	} else {
		log.Fatal("%v", ServeTLS(address, certPem, keyPem))
	}
}

// Local Variables:
// compile-cmd: "go test"
// End:
