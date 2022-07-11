package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianosela/rbac/api/model"
	"github.com/adrianosela/rbac/api/service/payloads"
	"github.com/adrianosela/rbac/utils/set"
	"github.com/gorilla/mux"
)

func (s *service) setPermissionEndpoints() {
	s.router.Methods(http.MethodPost).Path("/permission").Handler(s.auth(s.createPermissionHandler))
	s.router.Methods(http.MethodGet).Path("/permission/{name}").HandlerFunc(s.readPermissionHandler)
	s.router.Methods(http.MethodPatch).Path("/permission/{name}").Handler(s.auth(s.updatePermissionHandler))
	s.router.Methods(http.MethodPatch).Path("/permission/{name}/add").Handler(s.auth(s.addToPermissionHandler))         // add owners
	s.router.Methods(http.MethodPatch).Path("/permission/{name}/remove").Handler(s.auth(s.removeFromPermissionHandler)) // rm owners
	s.router.Methods(http.MethodDelete).Path("/permission/{name}").Handler(s.auth(s.deletePermissionHandler))
}

func (s *service) createPermissionHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	var pl *payloads.CreatePermissionRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request body is not a permission"))
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

	permission := &model.Permission{
		Name:        pl.Name,
		Description: pl.Description,
		Owners:      set.NewSet(pl.Owners...).Add(authenticatedUser).Slice(),
	}
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
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no permission name in request URL"))
		return
	}

	var pl *payloads.GenericUpdateDescriptionRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto an GenericUpdateDescriptionRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

	perm, err := s.store.ReadPermission(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read permission from storage"))
		return
	}
	if perm == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Permission \"%s\" does not exist!", name)))
		return
	}

	if !set.NewSet(perm.Owners...).Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Only the owners of a permission can modify the permission. User \"%s\" not in %v.", authenticatedUser, perm.Owners)))
		return
	}

	perm.Description = pl.Description
	if err := s.store.UpdatePermission(perm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update permission in storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Permission \"%s\" updated successfully!", name)))
	return
}

func (s *service) addToPermissionHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no permission name in request URL"))
		return
	}

	var pl *payloads.ModifyPermissionRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto an ModifyPermissionRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

	perm, err := s.store.ReadPermission(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read permission from storage"))
		return
	}
	if perm == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Permission \"%s\" does not exist!", name)))
		return
	}

	permOwners := set.NewSet(perm.Owners...)
	if !permOwners.Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Only the owners of a permission can modify the permission. User \"%s\" not in %v.", authenticatedUser, perm.Owners)))
		return
	}

	perm.Owners = permOwners.Add(pl.Owners...).Slice()
	if err := s.store.UpdatePermission(perm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update permission in storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Permission \"%s\" updated successfully!", name)))
	return
}

func (s *service) removeFromPermissionHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no permission name in request URL"))
		return
	}

	var pl *payloads.ModifyPermissionRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto an ModifyPermissionRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

	if set.NewSet(pl.Owners...).Has(authenticatedUser) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Removing yourself as an owner is not allowed"))
		return
	}

	perm, err := s.store.ReadPermission(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read permission from storage"))
		return
	}
	if perm == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Permission \"%s\" does not exist!", name)))
		return
	}

	permOwners := set.NewSet(perm.Owners...)
	if !permOwners.Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Only the owners of a permission can modify the permission. User \"%s\" not in %v.", authenticatedUser, perm.Owners)))
		return
	}

	perm.Owners = permOwners.Remove(pl.Owners...).Slice()
	if err := s.store.UpdatePermission(perm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update permission in storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Permission \"%s\" updated successfully!", name)))
	return
}

func (s *service) deletePermissionHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no permission name in request URL"))
		return
	}

	perm, err := s.store.ReadPermission(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read permission from storage"))
		return
	}
	if perm == nil { // (not in store already)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Permission \"%s\" deleted successfully!", name)))
		return
	}

	if !set.NewSet(perm.Owners...).Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("only the owners of a permission can delete the permission. User \"%s\" not in %v.", authenticatedUser, perm.Owners)))
		return
	}

	if len(perm.Roles) > 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("Permission \"%s\" is in use. Must first remove it from roles %v ", perm.Name, perm.Roles)))
		return
	}

	if err := s.store.DeletePermission(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete permission from storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Permission \"%s\" deleted successfully!", name)))
	return
}
