This simple email-sending service provides a wrapper/abstraction over other mail-sending services. Currently supported are SendGrid, MailGun, Mandrill, Amazon SES.

- go get github.com/erjoalgo/simple-mail-send-go
- Generate required cert.pem, key.pem files.

  - go run ${GOROOT}/src/crypto/tls/generate_cert.go --host localhost
- Start the server:

  - simple-mail-send-go --address localhost:6700 --keyPem key.pem --certPem cert.pem 
- Check server is running:

  - curl https://localhost:6700/ok -k
  - curl https://localhost:6700/credentials -k 
  - curl https://localhost:6700/sampleRequest -k 

- Send an example mail-send request:

::

	  curl -d @- -k https://localhost:6700/send <<END
	  {
		    "FromEmail": "me@example.com",
		    "ToEmails": ["you@example.com"],
		    "Subject":  "hola from simple-mail-send-go", 
		    "Credentials": {
			"mandrill": {"apikey":"<key>"},
			"mailgun": {"apikey":"<key>", "domain":"<domain>"},
			"sendgrid": {"user": "<user>", "passwd":  "<pass>"},
			"amazon": {"access-key-id":"<key>", "secret-access-key":"<key>"}
		    }
	}
	END



- The following fields are always required:

  - "FromEmail", a string
  - "ToEmails", a list of strings
  - "Credentials" a JSON map of service names to the map of required key-value pairs for each service. (See the "credential" endpoint for service-specific info)
