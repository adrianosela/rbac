package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/rbac/api/model"
	"github.com/adrianosela/rbac/utils/set"
	"github.com/gorilla/mux"
)

func (s *service) setRoleEndpoints() {
	s.router.Methods(http.MethodPost).Path("/role").HandlerFunc(s.createRoleHandler)
	s.router.Methods(http.MethodGet).Path("/role/{name}").HandlerFunc(s.readRoleHandler)
	s.router.Methods(http.MethodPatch).Path("/role/{name}").HandlerFunc(s.updateRoleHandler)
	s.router.Methods(http.MethodDelete).Path("/role/{name}").HandlerFunc(s.deleteRoleHandler)
}

func (s *service) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get authenticated user from ctx (for ownership)
	var role *model.Role
	if err := unmarshalRequestBody(r, &role); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body is not a role"))
		return
	}
	// TODO: add authenticated user to owners if not present
	user := "MOCK_AUTHENTICATED_USER"
	role.Owners = set.NewSet(role.Owners...).Add(user).Slice()
	// TODO: validate role has mandatory fields populated
	if err := s.store.CreateRole(role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to create new role in storage"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Role \"%s\" created successfuly!", role.Name)))
	return
}

func (s *service) readRoleHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no role name in request URL"))
		return
	}

	role, err := s.store.ReadRole(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read role from storage"))
		return
	}
	if role == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Role \"%s\" does not exist!", name)))
		return
	}

	roleBytes, err := json.Marshal(&role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(roleBytes)
	return
}

func (s *service) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
