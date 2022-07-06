package model

// Group represents a group
type Group struct {
	ID          string   `json:"id"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}
