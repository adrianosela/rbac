package service

import "net/http"

func (s *service) setRoleEndpoints() {
	s.router.Methods(http.MethodPost).Path("/role").HandlerFunc(s.createRoleHandler)
	s.router.Methods(http.MethodGet).Path("/role/{name}").HandlerFunc(s.readRoleHandler)
	s.router.Methods(http.MethodPatch).Path("/role/{name}").HandlerFunc(s.updateRoleHandler)
	s.router.Methods(http.MethodDelete).Path("/role/{name}").HandlerFunc(s.deleteRoleHandler)
}

func (s *service) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) readRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
