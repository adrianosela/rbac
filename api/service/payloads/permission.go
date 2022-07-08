package payloads

type CreatePermissionRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Owners      []string `json:"owners,omitempty"`
}
