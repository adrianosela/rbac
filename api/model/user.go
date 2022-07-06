package model

// User represents a user
type User struct {
	ID          string   `json:"id"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}
