package model

// Role represents a role
type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Assumers    []string `json:"assumers"`
	Owners      []string `json:"owners"`
}
