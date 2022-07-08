package storage

import (
	"fmt"

	"github.com/adrianosela/rbac/api/model"
	"github.com/adrianosela/rbac/utils/set"
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
	return nil, nil
}

// BulkReadPermissions retrieves a list of permission in storage
func (ms *MemoryStorage) BulkReadPermissions(names []string) ([]*model.Permission, error) {
	perms := []*model.Permission{}
	for _, name := range names {
		p, ok := ms.permissions[name]
		if !ok {
			return nil, fmt.Errorf("permission \"%s\" does not exist", name)
		}
		perms = append(perms, p)
	}
	return perms, nil
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
	return nil, nil
}

// BulkReadRoles retrieves a list of roles in storage
func (ms *MemoryStorage) BulkReadRoles(names []string) ([]*model.Role, error) {
	roles := []*model.Role{}
	for _, name := range names {
		r, ok := ms.roles[name]
		if !ok {
			return nil, fmt.Errorf("role \"%s\" does not exist", name)
		}
		roles = append(roles, r)
	}
	return roles, nil
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

// ReadUser retrieves a user in storage
func (ms *MemoryStorage) ReadUser(name string) (*model.User, error) {
	if r, ok := ms.users[name]; ok {
		return r, nil
	}
	return nil, nil
}

// AddRoleToUsers adds a role to the list of roles for users in storage.
// If the user does not exist, it is created
func (ms *MemoryStorage) AddRoleToUsers(role string, users []string) error {
	for _, user := range users {
		if u, ok := ms.users[user]; ok {
			u.Roles = set.NewSet(u.Roles...).Add(role).Slice()
		} else {
			ms.users[user] = &model.User{ID: user, Roles: []string{role}}
		}
	}
	return nil
}

// RemoveRoleFromUsers removes a role from the list of roles for users in storage.
func (ms *MemoryStorage) RemoveRoleFromUsers(role string, users []string) error {
	for _, user := range users {
		if u, ok := ms.users[user]; ok {
			u.Roles = set.NewSet(u.Roles...).Remove(role).Slice()
		}
	}
	return nil
}

// UpdateUser updates a user in storage
func (ms *MemoryStorage) UpdateUser(u *model.User) error {
	ms.users[u.ID] = u
	return nil
}

// DeleteUser deletes a user in storage
func (ms *MemoryStorage) DeleteUser(id string) error {
	delete(ms.users, id)
	return nil
}

// ReadGroup retrieves a group in storage
func (ms *MemoryStorage) ReadGroup(id string) (*model.Group, error) {
	if r, ok := ms.groups[id]; ok {
		return r, nil
	}
	return nil, nil
}

// ReadGroups retrieves a list of groups in storage
// NOTE: behavior for not found groups differs than from not found in bulk roles/perms
func (ms *MemoryStorage) ReadGroups(names []string) ([]*model.Group, error) {
	groups := []*model.Group{}
	for _, name := range names {
		g, ok := ms.groups[name]
		if !ok {
			continue
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// AddRoleToGroups adds a role to the list of roles for groups in storage.
// If the group does not exist, it is created
func (ms *MemoryStorage) AddRoleToGroups(role string, groups []string) error {
	for _, group := range groups {
		if g, ok := ms.groups[group]; ok {
			g.Roles = set.NewSet(g.Roles...).Add(role).Slice()
		} else {
			ms.groups[group] = &model.Group{ID: group, Roles: []string{role}}
		}
	}
	return nil
}

// RemoveRoleFromGroups removes a role from the list of roles for groups in storage.
func (ms *MemoryStorage) RemoveRoleFromGroups(role string, groups []string) error {
	for _, group := range groups {
		if g, ok := ms.groups[group]; ok {
			g.Roles = set.NewSet(g.Roles...).Remove(role).Slice()
		}
	}
	return nil
}

// UpdateGroup updates a group in storage
func (ms *MemoryStorage) UpdateGroup(g *model.Group) error {
	ms.groups[g.ID] = g
	return nil
}

// DeleteGroup deletes a group in storage
func (ms *MemoryStorage) DeleteGroup(id string) error {
	delete(ms.groups, id)
	return nil
}
