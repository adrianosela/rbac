package groups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// OktaSource is an implementation of the Source interface that
// leverages Okta as the source-of-truth for group memberships.
type OktaSource struct {
	orgOktaDomain string // e.g. "your_company.oktapreview" or "your_company.okta"
	apiToken      string
	httpClient    *http.Client
}

// NewOktaSource returns a new OktaSource.
func NewOktaSource(orgOktaDomain, apiToken string) *OktaSource {
	return &OktaSource{
		orgOktaDomain: orgOktaDomain,
		apiToken:      apiToken,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GetForUser returns the groups a given user (id or login) is a member of.
// https://developer.okta.com/docs/reference/api/users/#get-user-s-groups
func (os *OktaSource) GetForUser(id string) ([]string, error) {
	url := fmt.Sprintf("https://%s.com/api/v1/users/%s/groups", os.orgOktaDomain, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to build http request: %s", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("SSWS %s", os.apiToken))

	resp, err := os.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make http request: %s", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return []string{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got a non 200 HTTP status code: %d", resp.StatusCode)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read http response body: %s", err)
	}
	defer resp.Body.Close()

	var oktaGroupsResponse []struct {
		Profile struct {
			Name string `json:"name"`
		} `json:"profile"`
	}

	if err = json.Unmarshal(respBodyBytes, &oktaGroupsResponse); err != nil {
		return nil, fmt.Errorf("Failed to decode HTTP response body: %s", err)

	}

	names := []string{}
	for _, group := range oktaGroupsResponse {
		names = append(names, group.Profile.Name)
	}

	return names, nil
}
