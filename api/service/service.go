package service

import (
	"net/http"

	"github.com/adrianosela/rbac/api/groups"
	"github.com/adrianosela/rbac/api/storage"
	"github.com/gorilla/mux"
)

var (
	MOCK_AUTHENTICATED_USER = "adriano" // FIXME
)

// Config represents configuration for the service
type Config struct {
	OktaOrgDomain string
	OktaAPIToken  string
}

type service struct {
	router *mux.Router
	store  storage.Storage
	groups groups.Source
}

// New returns the handler for a new service
func New(c Config) (http.Handler, error) {

	svc := &service{
		router: mux.NewRouter(),
		store:  storage.NewMemoryStorage(), // FIXME: use remote storage
		groups: groups.NewOktaSource(c.OktaOrgDomain, c.OktaAPIToken),
	}

	svc.setDebugEndpoints()
	svc.setPermissionEndpoints()
	svc.setRoleEndpoints()
	svc.setUserEndpoints()

	return svc.router, nil
}
