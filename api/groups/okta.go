package groups

import (
	"fmt"
)

// OktaSource is an implementation of the Source interface
// that leverages Okta as the source-of-truth for group memberships
type OktaSource struct {
	// TODO
}

// NewOktaSource returns a new OktaSource
func NewOktaSource( /* TODO */ ) *OktaSource {
	return &OktaSource{ /* TODO */ }
}

// GetForUser returns the groups a given user is a member of
func (os *OktaSource) GetForUser(id string) ([]string, error) {
	// TODO
	return nil, fmt.Errorf("OktaSource not implemented")
}
