package payloads

type GetUserPermissionsResponse struct {
	Persmissions []string `json:"permissions"`
}
