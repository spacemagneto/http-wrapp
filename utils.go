package client

// CopySet creates a shallow copy of a set (represented as a map with empty struct values).
// It is used to prevent side effects when the original set needs to be modified
// independently of the copy.
//
// Constraints: The key type K must satisfy the 'comparable' constraint,
// which is required for all map keys in Go.
func CopySet[K comparable](src map[K]struct{}) map[K]struct{} {
	dst := make(map[K]struct{}, len(src))
	for k := range src {
		dst[k] = struct{}{}
	}

	return dst
}
