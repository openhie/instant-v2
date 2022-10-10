package core

import (
	"os"
	"testing"

	"github.com/luno/jettison/jtest"
)

func TestGenerateConfigFile(t *testing.T) {
	type cases struct {
		config      Config
		expectError error
	}

	testCases := []cases{
		{
			config: Config{
				Image:         "implementation-image",
				ProjectName:   "test-project",
				PlatformImage: "jembi/go-cli-test-image:latest",
			},
		},
		{
			config: Config{
				ProjectName:   "test-project",
				PlatformImage: "jembi/go-cli-test-image:latest",
			},
			expectError: ErrInvalidConfig,
		},
		{
			config: Config{
				Image:         "implementation-image",
				PlatformImage: "jembi/go-cli-test-image:latest",
			},
			expectError: ErrInvalidConfig,
		},
		{
			config: Config{
				Image:       "implementation-image",
				ProjectName: "test-project",
			},
			expectError: ErrInvalidConfig,
		},
	}

	for _, tc := range testCases {
		defer os.Remove("config.yaml")
		err := GenerateConfigFile(&tc.config)
		jtest.Require(t, tc.expectError, err)
	}
}
