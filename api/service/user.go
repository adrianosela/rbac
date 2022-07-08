package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/rbac/api/service/payloads"
	"github.com/adrianosela/rbac/utils/set"
	"github.com/gorilla/mux"
)

func (s *service) setUserEndpoints() {
	s.router.Methods(http.MethodGet).Path("/user/{name}").HandlerFunc(s.getUserPermissionsHandler)
}

func (s *service) getUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no role name in request URL"))
		return
	}

	groups, err := s.groups.GetForUser(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get groups for user: %s", err))) // FIXME: do not expose internals
		return
	}

	roles := set.NewSet()

	// collect roles tied to groups
	gs, err := s.store.ReadGroups(groups)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to bulk-get groups from store: %s", err))) // FIXME: do not expose internals
		return
	}
	for _, group := range gs {
		roles.Add(group.Roles...)
	}

	// collect roles tied to user
	user, err := s.store.ReadUser(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get user from store: %s", err))) // FIXME: do not expose internals
		return
	}
	if user != nil {
		roles.Add(user.Roles...)
	}

	rs, err := s.store.BulkReadRoles(roles.Slice())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to bulk-get roles from store: %s", err))) // FIXME: do not expose internals
		return
	}

	perms := set.NewSet()
	for _, role := range rs {
		perms.Add(role.Permissions...)
	}

	respBytes, err := json.Marshal(&payloads.GetUserPermissionsResponse{Persmissions: perms.Slice()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to marshal response: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return
}
