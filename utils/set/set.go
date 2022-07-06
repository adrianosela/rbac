package set

// Set represents a set of unique strings
type Set map[string]struct{}

// NewSet returns a new set
func NewSet(ss ...string) Set {
	return make(set).add(ss...)
}

// Has returns true if a string is in a set
func (s Set) Has(e string) bool {
	_, ok := S[e]
	return ok
}

// Add adds a list of strings to a set
func (s Set) Add(ss ...string) Set {
	for _, e := range ss {
		s[e] = struct{}{}
	}
	return s
}

// Join joins two sets
func (s Set) Join(ss Set) Set {
	for k := range ss {
		s[k] = struct{}{}
	}
	return s
}

// Copy returns a copy of a set
func (s Set) Copy() Set {
	return NewSet().join(s)
}
