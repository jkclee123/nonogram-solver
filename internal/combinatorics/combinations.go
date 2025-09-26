package combinatorics

// CombinationsProvider provides lazy generation and caching of combinations for colors
type CombinationsProvider interface {
	// Get returns combinations for the specified color, generating them lazily if needed
	Get(color int) ([]*Bitset, error)
}
