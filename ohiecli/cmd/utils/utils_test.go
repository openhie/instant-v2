package utils

import (
	"fmt"
	"testing"
)

func Test_sliceContains(t *testing.T) {
	testCases := []struct {
		slice   []string
		element string
		result  bool
		name    string
	}{
		{
			name:    "SliceContain test - should return true when slice contains element",
			slice:   []string{"Optimus Prime", "Iron Hyde"},
			element: "Optimus Prime",
			result:  true,
		},
		{
			name:    "SliceContain test - should return false when slice does not contain element",
			slice:   []string{"Optimus Prime", "Iron Hyde"},
			element: "Megatron",
			result:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ans := SliceContains(tt.slice, tt.element)

			if ans != tt.result {
				t.Fatal("SliceContains should return " + fmt.Sprintf("%t", tt.result) + " but returned " + fmt.Sprintf("%t", ans))
			}
			t.Log(tt.name + " passed!")
		})
	}
}
