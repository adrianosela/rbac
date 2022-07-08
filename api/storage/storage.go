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

	CreateRole(*model.Role) error
	ReadRole(string) (*model.Role, error)
	BulkReadRoles([]string) ([]*model.Role, error)
	UpdateRole(*model.Role) error
	DeleteRole(string) error

	ReadUser(string) (*model.User, error)
	AddRoleToUsers(string, []string) error
	RemoveRoleFromUsers(string, []string) error

	ReadGroups([]string) ([]*model.Group, error)
	AddRoleToGroups(string, []string) error
	RemoveRoleFromGroups(string, []string) error
}
