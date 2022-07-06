package model

// Permission represents a permission
type Permission struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owners      []string `json:"owners"`
}
