package storage

import (
	"github.com/adrianosela/rbac/api/model"
)

// Storage represents the storage needs of the api
type Storage interface {
	CreatePermission(*model.Permission) error
	ReadPermission(string) (*model.Permission, error)
	UpdatePermission(*model.Permission) error
	DeletePermission(string) error

	CreateRole(*model.Role) error
	ReadRole(string) (*model.Role, error)
	UpdateRole(*model.Role) error
	DeleteRole(string) error

	CreateUser(*model.User) error
	ReadUser(string) (*model.User, error)
	UpdateUser(*model.User) error
	DeleteUser(string) error

	CreateGroup(*model.Group) error
	ReadGroup(string) (*model.Group, error)
	UpdateGroup(*model.Group) error
	DeleteGroup(string) error
}
