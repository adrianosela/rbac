package storage

import (
	"fmt"

	"github.com/adrianosela/rbac/api/model"
)

// MemoryStorage is an in-memory implementation of the Storage interface
type MemoryStorage struct {
	permissions map[string]*model.Permission
	roles       map[string]*model.Role
	users       map[string]*model.User
	groups      map[string]*model.Group
}

// NewMemoryStorage returns a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	ms := &MemoryStorage{
		permissions: make(map[string]*model.Permission),
		roles:       make(map[string]*model.Role),
		users:       make(map[string]*model.User),
		groups:      make(map[string]*model.Group),
	}
	return ms
}

// CreatePermission creates a new permission in storage
func (ms *MemoryStorage) CreatePermission(p *model.Permission) error {
	if _, ok := ms.permissions[p.Name]; ok {
		return fmt.Errorf("permission \"%s\" already exists", p.Name)
	}
	ms.permissions[p.Name] = p
	return nil
}

// ReadPermission retrieves a permission in storage
func (ms *MemoryStorage) ReadPermission(name string) (*model.Permission, error) {
	if p, ok := ms.permissions[name]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("permission \"%s\" does not exist", name)
}

// UpdatePermission updates a permission in storage
func (ms *MemoryStorage) UpdatePermission(p *model.Permission) error {
	if _, ok := ms.permissions[p.Name]; !ok {
		return fmt.Errorf("permission \"%s\" does not exist", p.Name)
	}
	ms.permissions[p.Name] = p
	return nil
}

// DeletePermission deletes a permission in storage
func (ms *MemoryStorage) DeletePermission(name string) error {
	delete(ms.permissions, name)
	return nil
}

// CreateRole creates a new role in storage
func (ms *MemoryStorage) CreateRole(r *model.Role) error {
	if _, ok := ms.roles[r.Name]; ok {
		return fmt.Errorf("role \"%s\" already exists", r.Name)
	}
	ms.roles[r.Name] = r
	return nil
}

// ReadRole retrieves a role in storage
func (ms *MemoryStorage) ReadRole(name string) (*model.Role, error) {
	if r, ok := ms.roles[name]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("role \"%s\" does not exist", name)
}

// UpdateRole updates a role in storage
func (ms *MemoryStorage) UpdateRole(r *model.Role) error {
	if _, ok := ms.roles[r.Name]; !ok {
		return fmt.Errorf("role \"%s\" does not exist", r.Name)
	}
	ms.roles[r.Name] = r
	return nil
}

// DeleteRole deletes a role in storage
func (ms *MemoryStorage) DeleteRole(name string) error {
	delete(ms.roles, name)
	return nil
}

// CreateUser creates a new user in storage
func (ms *MemoryStorage) CreateUser(u *model.User) error {
	if _, ok := ms.users[u.ID]; ok {
		return fmt.Errorf("user \"%s\" already exists", u.ID)
	}
	ms.users[u.ID] = u
	return nil
}

// ReadUser retrieves a user in storage
func (ms *MemoryStorage) ReadUser(name string) (*model.User, error) {
	if r, ok := ms.users[name]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("user \"%s\" does not exist", name)
}

// UpdateUser updates a user in storage
func (ms *MemoryStorage) UpdateUser(u *model.User) error {
	if _, ok := ms.users[u.ID]; !ok {
		return fmt.Errorf("user \"%s\" does not exist", u.ID)
	}
	ms.users[u.ID] = u
	return nil
}

// DeleteUser deletes a user in storage
func (ms *MemoryStorage) DeleteUser(id string) error {
	delete(ms.users, id)
	return nil
}

// CreateGroup creates a new group in storage
func (ms *MemoryStorage) CreateGroup(g *model.Group) error {
	if _, ok := ms.groups[g.ID]; ok {
		return fmt.Errorf("group \"%s\" already exists", g.ID)
	}
	ms.groups[g.ID] = g
	return nil
}

// ReadGroup retrieves a group in storage
func (ms *MemoryStorage) ReadGroup(id string) (*model.Group, error) {
	if r, ok := ms.groups[id]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("group \"%s\" does not exist", id)
}

// UpdateGroup updates a group in storage
func (ms *MemoryStorage) UpdateGroup(g *model.Group) error {
	if _, ok := ms.groups[g.ID]; !ok {
		return fmt.Errorf("group \"%s\" does not exist", g.ID)
	}
	ms.groups[g.ID] = g
	return nil
}

// DeleteGroup deletes a group in storage
func (ms *MemoryStorage) DeleteGroup(id string) error {
	delete(ms.groups, id)
	return nil
}
