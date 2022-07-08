package model

// Role represents a role
type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Users       []string `json:"users"`
	Groups      []string `json:"groups"`
	Owners      []string `json:"owners"`
}
