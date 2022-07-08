package model

// User represents a user
type User struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
	// Permissions []string `json:"permissions"`
}
