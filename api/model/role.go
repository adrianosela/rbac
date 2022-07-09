package model

// Role represents a role
type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owners      []string `json:"owners"`
	Users       []string `json:"users"`
	Groups      []string `json:"groups"`
	Permissions []string `json:"permissions"`
}
