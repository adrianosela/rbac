package service

import "net/http"

func (s *service) setPermissionEndpoints() {
	s.router.Methods(http.MethodPost).Path("/permission").HandlerFunc(s.createPermissionHandler)
	s.router.Methods(http.MethodGet).Path("/permission/{name}").HandlerFunc(s.readPermissionHandler)
	s.router.Methods(http.MethodPatch).Path("/permission/{name}").HandlerFunc(s.updatePermissionHandler)
	s.router.Methods(http.MethodDelete).Path("/permission/{name}").HandlerFunc(s.deletePermissionHandler)
}

func (s *service) createPermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) readPermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) updatePermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) deletePermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
