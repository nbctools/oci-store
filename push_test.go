package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string
		want     string
	}{
		{
			name:     "existing env var",
			key:      "TEST_VAR",
			setValue: "test_value",
			want:     "test_value",
		},
		{
			name:     "non-existing env var",
			key:      "NON_EXISTING_VAR",
			setValue: "",
			want:     "",
		},
		{
			name:     "env var with spaces",
			key:      "TEST_SPACES",
			setValue: "  value with spaces  ",
			want:     "value with spaces",
		},
		{
			name:     "empty env var",
			key:      "TEST_EMPTY",
			setValue: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			got := getEnv(tt.key)
			if got != tt.want {
				t.Errorf("getEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}
