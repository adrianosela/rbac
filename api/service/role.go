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

func (s *service) setRoleEndpoints() {
	s.router.Methods(http.MethodPost).Path("/role").Handler(s.auth(s.createRoleHandler))

	s.router.Methods(http.MethodGet).Path("/role/{name}").HandlerFunc(s.readRoleHandler)

	s.router.Methods(http.MethodPatch).Path("/role/{name}").Handler(s.auth(s.updateRoleHandler))            // modify description
	s.router.Methods(http.MethodPatch).Path("/role/{name}/add").Handler(s.auth(s.addToRoleHandler))         // add permissions, assumers, or owners
	s.router.Methods(http.MethodPatch).Path("/role/{name}/remove").Handler(s.auth(s.removeFromRoleHandler)) // rm permissions, assumers, or owners

	s.router.Methods(http.MethodDelete).Path("/role/{name}").Handler(s.auth(s.deleteRoleHandler))
}

func (s *service) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	var pl *payloads.CreateRoleRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto a CreateRoleRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

	role := &model.Role{
		Name:        pl.Name,
		Description: pl.Description,
		Permissions: set.NewSet(pl.Permissions...).Slice(),
		Users:       set.NewSet(pl.Users...).Slice(),
		Groups:      set.NewSet(pl.Groups...).Slice(),
		Owners:      set.NewSet(pl.Owners...).Add(authenticatedUser).Slice(),
	}

	perms, err := s.store.BulkReadPermissions(pl.Permissions)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // FIXME: do not expose internals
		return
	}

	for _, perm := range perms {
		if !set.NewSet(perm.Owners...).Has(authenticatedUser) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Only the owners of a permission can add it to a role. User \"%s\" not in %v.", authenticatedUser, perm.Owners)))
			return
		}
	}

	if err := s.store.CreateRole(role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to create new role in storage"))
		return
	}

	// FIXME: move three below to eventual consistence model
	if err := s.store.AddRoleToPermissions(pl.Name, pl.Permissions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add role to permissions in storage"))
		return
	}
	if err := s.store.AddRoleToUsers(pl.Name, pl.Users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add role to users in storage"))
		return
	}
	if err := s.store.AddRoleToGroups(pl.Name, pl.Groups); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add role to groups in storage"))
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
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no role name in request URL"))
		return
	}

	var pl *payloads.GenericUpdateDescriptionRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto an UpdateRoleRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

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

	if !set.NewSet(role.Owners...).Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Only the owners of a role can modify the role. User \"%s\" not in %v.", authenticatedUser, role.Owners)))
		return
	}

	role.Description = pl.Description
	if err := s.store.UpdateRole(role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update role in storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Role \"%s\" updated successfully!", name)))
	return
}

func (s *service) addToRoleHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no role name in request URL"))
		return
	}

	var pl *payloads.ModifyRoleRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto an ModifyRoleRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

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

	roleOwners := set.NewSet(role.Owners...)
	if !roleOwners.Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Only the owners of a role can modify the role. User \"%s\" not in %v.", authenticatedUser, role.Owners)))
		return
	}

	perms, err := s.store.BulkReadPermissions(pl.Permissions)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // FIXME: do not expose internals
		return
	}
	for _, perm := range perms {
		if !set.NewSet(perm.Owners...).Has(authenticatedUser) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Only the owners of a permission can add it to a role. User \"%s\" not in %v.", authenticatedUser, perm.Owners)))
			return
		}
	}

	role.Owners = roleOwners.Add(pl.Owners...).Slice()
	role.Users = set.NewSet(role.Users...).Add(pl.Users...).Slice()
	role.Groups = set.NewSet(role.Groups...).Add(pl.Groups...).Slice()
	role.Permissions = set.NewSet(role.Permissions...).Add(pl.Permissions...).Slice()

	if err := s.store.UpdateRole(role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update role in storage"))
		return
	}

	// FIXME: move three below to eventual consistence model
	if err := s.store.AddRoleToPermissions(name, pl.Permissions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add role to permissions in storage"))
		return
	}
	if err := s.store.AddRoleToUsers(name, pl.Users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add role to users in storage"))
		return
	}
	if err := s.store.AddRoleToGroups(name, pl.Groups); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to add role to groups in storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Role \"%s\" updated successfully!", name)))
	return
}

func (s *service) removeFromRoleHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

	name := mux.Vars(r)["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no role name in request URL"))
		return
	}

	var pl *payloads.ModifyRoleRequest
	if err := unmarshalRequestBody(r, &pl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not decode request body onto an ModifyRoleRequest")) // FIXME: don't expose internals
		return
	}

	// TODO: validate payload (e.g. for required fields, length limits, etc)

	if set.NewSet(pl.Owners...).Has(authenticatedUser) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Removing yourself as an owner is not allowed"))
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

	owners := set.NewSet(role.Owners...)
	if !owners.Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Only the owners of a role can modify the role. User \"%s\" not in %v.", authenticatedUser, role.Owners)))
		return
	}

	role.Owners = owners.Remove(pl.Owners...).Slice()
	role.Users = set.NewSet(role.Users...).Remove(pl.Users...).Slice()
	role.Groups = set.NewSet(role.Groups...).Remove(pl.Groups...).Slice()
	role.Permissions = set.NewSet(role.Permissions...).Remove(pl.Permissions...).Slice()

	if err := s.store.UpdateRole(role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update role in storage"))
		return
	}

	// FIXME: move three below to eventual consistence model
	if err := s.store.RemoveRoleFromPermissions(name, pl.Permissions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to remove role from permissions in storage"))
		return
	}
	if err := s.store.RemoveRoleFromUsers(name, pl.Users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to remove role from users in storage"))
		return
	}
	if err := s.store.RemoveRoleFromGroups(name, pl.Groups); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to remove role from groups in storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Role \"%s\" updated successfully!", name)))
	return
}

func (s *service) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := getAuthenticatedUser(r)

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

	if role == nil { // (not in store already)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Role \"%s\" deleted successfully!", name)))
		return
	}

	if !set.NewSet(role.Owners...).Has(authenticatedUser) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("only the owners of a role can delete the role. User \"%s\" not in %v.", authenticatedUser, role.Owners)))
		return
	}

	// FIXME: move three below to eventual consistence model
	if err := s.store.RemoveRoleFromPermissions(name, role.Permissions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to remove role from permissions in storage"))
		return
	}
	if err := s.store.RemoveRoleFromUsers(name, role.Users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to remove role from users in storage"))
		return
	}
	if err := s.store.RemoveRoleFromGroups(name, role.Groups); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to remove role from groups in storage"))
		return
	}

	if err := s.store.DeleteRole(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete role from storage"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Role \"%s\" deleted successfully!", name)))
	return
}
