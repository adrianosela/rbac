package storage

import (
	"github.com/adrianosela/rbac/api/model"
)

// Storage represents the storage needs of the api
type Storage interface {
	CreatePermission(*model.Permission) error
	ReadPermission(string) (*model.Permission, error)
	BulkReadPermissions([]string) ([]*model.Permission, error)
	UpdatePermission(*model.Permission) error
	DeletePermission(string) error
	AddRoleToPermissions(string, []string) error      // FIXME: move to eventual consistence
	RemoveRoleFromPermissions(string, []string) error // FIXME: move to eventual consistence

	CreateRole(*model.Role) error
	ReadRole(string) (*model.Role, error)
	BulkReadRoles([]string) ([]*model.Role, error)
	UpdateRole(*model.Role) error
	DeleteRole(string) error

	ReadUser(string) (*model.User, error)
	AddRoleToUsers(string, []string) error      // FIXME: move to eventual consistence
	RemoveRoleFromUsers(string, []string) error // FIXME: move to eventual consistence

	ReadGroups([]string) ([]*model.Group, error)
	AddRoleToGroups(string, []string) error      // FIXME: move to eventual consistence
	RemoveRoleFromGroups(string, []string) error // FIXME: move to eventual consistence
}
