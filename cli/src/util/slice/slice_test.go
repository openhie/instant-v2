package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SliceContains(t *testing.T) {
	type testCase struct {
		slice      interface{}
		typeOf     string
		matchValue interface{}
		wantMatch  bool
	}

	testCases := []testCase{
		{
			slice:      []string{"one", "two", "three"},
			typeOf:     "string",
			matchValue: "two",
			wantMatch:  true,
		},
		{
			slice:      []string{"one", "two", "three"},
			typeOf:     "string",
			matchValue: "four",
			wantMatch:  false,
		},
		{
			slice:      []int{1, 2, 3},
			typeOf:     "int",
			matchValue: 2,
			wantMatch:  true,
		},
		{
			slice:      []int{1, 2, 3},
			typeOf:     "int",
			matchValue: 4,
			wantMatch:  false,
		},
		{
			slice:      []int{},
			typeOf:     "int",
			matchValue: 0,
			wantMatch:  false,
		},
		{
			slice:      []string{},
			typeOf:     "string",
			matchValue: "",
			wantMatch:  false,
		},
	}

	for _, tc := range testCases {
		switch tc.typeOf {
		case "string":
			require.Equal(t, SliceContains(tc.slice.([]string), tc.matchValue.(string)), tc.wantMatch)

		case "int":
			require.Equal(t, SliceContains(tc.slice.([]int), tc.matchValue.(int)), tc.wantMatch)

		}
	}
}
