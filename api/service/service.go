package service

import (
	"net/http"

	"github.com/adrianosela/rbac/api/storage"
	"github.com/gorilla/mux"
)

type service struct {
	router *mux.Router
	store  storage.Storage
}

// New returns the handler for a new service
func New() (http.Handler, error) {
	rtr := mux.NewRouter()

	svc := &service{
		router: rtr,
		store:  storage.NewMemoryStorage(),
	}

	svc.setDebugEndpoints()
	svc.setPermissionEndpoints()
	svc.setRoleEndpoints()
	svc.setUserEndpoints()

	return svc.router, nil
}
