package slice

import "strings"

func SliceContains[Type comparable](slice []Type, element Type) bool {
	for _, s := range slice {
		if element == s {
			return true
		}
	}

	return false
}

func AppendUniqueToMapFromSlice(m map[string]string, sl []string) map[string]string {
	for _, s := range sl {
		splitEnvVar := strings.SplitAfter(s, "=")

		if _, ok := m[splitEnvVar[0]]; !ok {
			m[splitEnvVar[0]] = splitEnvVar[1]
		}
	}

	return m
}
