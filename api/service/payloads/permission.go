package payloads

type CreatePermissionRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owners      []string `json:"owners,omitempty"`
}

type ModifyPermissionRequest struct {
	Owners []string `json:"owners,omitempty"`
}
