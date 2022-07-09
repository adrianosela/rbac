package payloads

type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
	Users       []string `json:"users,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Owners      []string `json:"owners,omitempty"`
}

type ModifyRoleRequest struct {
	Permissions []string `json:"permissions,omitempty"`
	Users       []string `json:"users,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Owners      []string `json:"owners,omitempty"`
}
