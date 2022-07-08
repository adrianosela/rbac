package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adrianosela/rbac/api/service"
)

func main() {
	config := service.Config{
		OktaOrgDomain: os.Getenv("OKTA_ORG_DOMAIN"),
		OktaAPIToken:  os.Getenv("OKTA_API_TOKEN"),
	}

	handler, err := service.New(config)
	if err != nil {
		log.Fatalf("Failed to initialize service: %s", err)
	}

	if err = http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to serve HTTP: %s", err)
	}
}
