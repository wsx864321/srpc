package util

// InSliceString ...
func InSliceString(s string, strs []string) bool {
	for _, val := range strs {
		if val == s {
			return true
		}
	}

	return false
}
