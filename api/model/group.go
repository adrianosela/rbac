package model

// Group represents a group
type Group struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
	// Permissions []string `json:"permissions"`
}
