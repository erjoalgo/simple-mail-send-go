This simple email-sending service provides a wrapper/abstraction over other mail-sending services. Currently supported are SendGrid, MailGun, Mandrill, Amazon SES.

- go get https://github.com/erjoalgo/simple-mail-send-go.git
- Generate required cert.pem, key.pem files.

  - go run ${GOROOT}/src/crypto/tls/generate_cert.go --host localhost
- Start the server:

  - simple-mail-send-go --address localhost:6700 --keyPem key.pem --certPem cert.pem 
- Check server is running:

  - curl https://localhost:6700/ok -k 

- Send an example mail-send request:

.. code:: python
	  curl -d @- -k https://localhost:6700/send <<END
	  {
		    "FromEmail": "me@example.com",
		    "ToEmails": ["you@example.com"],
		    "Subject":  "hola from simple-mail-send-go", 
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
	}
END
