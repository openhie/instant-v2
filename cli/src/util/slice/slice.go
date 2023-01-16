package slice

func SliceContains[Type comparable](slice []Type, element Type) bool {
	for _, s := range slice {
		if element == s {
			return true
		}
	}

	return false
}
