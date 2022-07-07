package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/rbac/api/model"
	"github.com/gorilla/mux"
)

func (s *service) setPermissionEndpoints() {
	s.router.Methods(http.MethodPost).Path("/permission").HandlerFunc(s.createPermissionHandler)
	s.router.Methods(http.MethodGet).Path("/permission/{name}").HandlerFunc(s.readPermissionHandler)
	s.router.Methods(http.MethodPatch).Path("/permission/{name}").HandlerFunc(s.updatePermissionHandler)
	s.router.Methods(http.MethodDelete).Path("/permission/{name}").HandlerFunc(s.deletePermissionHandler)
}

func (s *service) createPermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get authenticated user from ctx (for ownership)
	var permission *model.Permission
	if err := unmarshalRequestBody(r, &permission); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body is not a permission"))
		return
	}
	// TODO: add authenticated user to owners if not present
	// TODO: validate role has mandatory fields populated
	if err := s.store.CreatePermission(permission); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to create new permission in storage"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Permission \"%s\" created successfuly!", permission.Name)))
	return
}

func (s *service) readPermissionHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no permission name in request URL"))
		return
	}

	permission, err := s.store.ReadPermission(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read permission from storage"))
		return
	}
	if permission == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Permission \"%s\" does not exist!", name)))
		return
	}

	permissionBytes, err := json.Marshal(&permission)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(permissionBytes)
	return
}

func (s *service) updatePermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *service) deletePermissionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
