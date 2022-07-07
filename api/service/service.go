package service

import (
	"net/http"

	"github.com/adrianosela/rbac/api/groups"
	"github.com/adrianosela/rbac/api/storage"
	"github.com/gorilla/mux"
)

type service struct {
	router *mux.Router
	store  storage.Storage
	groups groups.Source
}

var mockGroupMemberships = map[string][]string{
	"adriano@adrianosela.com": {"SIRT", "IT", "Engineering", "TeamAvocados"},
}

// New returns the handler for a new service
func New() (http.Handler, error) {
	rtr := mux.NewRouter()

	svc := &service{
		router: rtr,
		store:  storage.NewMemoryStorage(),                   // FIXME: use remote storage
		groups: groups.NewMemorySource(mockGroupMemberships), // FIXME: use remote groups source
	}

	svc.setDebugEndpoints()
	svc.setPermissionEndpoints()
	svc.setRoleEndpoints()
	svc.setUserEndpoints()

	return svc.router, nil
}
