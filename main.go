package main

import (
	"log"
	"net/http"
	"os"

	"codec-cors-credentials/platform/authenticator"
	"codec-cors-credentials/platform/router"
)

func main() {
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	bindAddress, ok := os.LookupEnv("CODEC_BIND_ADDRESS")
	if !ok {
		bindAddress = ":3000"
	}

	log.Print("Server listening on ", bindAddress)
	if err := http.ListenAndServe(bindAddress, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}

	// if err := http.ListenAndServeTLS(bindAddress, "certs/tls.crt", "certs/tls.key", rtr); err != nil {
	// 	log.Fatalf("There was an error with the http server: %v", err)
	// }
}
