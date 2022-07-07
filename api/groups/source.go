package groups

// Source represents the functionality of a groups source
type Source interface {
	GetForUser(string) ([]string, error)
}
