package service

import (
	"net/http"
)

func (s *service) setDebugEndpoints() {
	s.router.Methods(http.MethodGet).Path("/healthcheck").HandlerFunc(s.healthcheckHandler)
}

func (s *service) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm alive!"))
	return
}
