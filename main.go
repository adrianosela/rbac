package main

import (
	"log"
	"net/http"

	"github.com/adrianosela/rbac/api/service"
)

func main() {
	handler, err := service.New()
	if err != nil {
		log.Fatalf("Failed to initialize service: %s", err)
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to serve HTTP: %s", err)
	}
}
