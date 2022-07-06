package service

import "net/http"

func (s *service) setUserEndpoints() {
	s.router.Methods(http.MethodGet).Path("/user/{name}").HandlerFunc(s.getUserHandler)
}

func (s *service) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
