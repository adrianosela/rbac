package groups

import (
	"fmt"
)

// MemorySource is an in-memory implementation of the Source interface
type MemorySource struct {
	groups map[string][]string
}

// NewMemorySource returns a new MemorySource
func NewMemorySource(groups map[string][]string) *MemorySource {
	return &MemorySource{groups: groups}
}

// GetForUser returns the groups a given user is a member of
func (ms *MemorySource) GetForUser(id string) ([]string, error) {
	gm, ok := ms.groups[id]
	if !ok {
		return nil, fmt.Errorf("User \"%s\" not found", id)
	}
	return gm, nil
}
