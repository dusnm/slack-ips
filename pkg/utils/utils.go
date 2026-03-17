package utils

// MergeMaps
//
// Merges N maps into a new map.
// Duplicate keys will hold the value of the last map in the set.
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	totalLength := 0
	for _, m := range maps {
		totalLength += len(m)
	}

	merged := make(map[K]V, totalLength)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}
