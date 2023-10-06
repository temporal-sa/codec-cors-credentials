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

	bindAddress := getEnv("CODEC_BIND_ADDRESS", "localhost:3000")
	enableTLS := getEnv("CODEC_ENABLE_TLS", "true")

	if enableTLS == "true" {
		log.Printf("Server listening on https://%v\n", bindAddress)
		err = http.ListenAndServeTLS(bindAddress, "certs/tls.crt", "certs/tls.key", rtr)
	} else {
		log.Printf("Server listening on http://%v\n", bindAddress)
		err = http.ListenAndServe(bindAddress, rtr)
	}
	if err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
