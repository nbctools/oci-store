package main

import (
	"os"
	"testing"
)

func TestValidateGCSConfig(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		wantErr  bool
	}{
		{
			name: "valid config with keyfile flag",
			setup: func() {
				gcsKeyfile = "my-project.json"
			},
			teardown: func() {
				gcsKeyfile = ""
			},
			wantErr: false,
		},
		{
			name: "valid config with env var",
			setup: func() {
				gcsKeyfile = ""
				os.Setenv("GOOGLE_CLOUD_PROJECT", "env-project")
			},
			teardown: func() {
				gcsKeyfile = ""
				os.Unsetenv("GOOGLE_CLOUD_PROJECT")
			},
			wantErr: false,
		},
		{
			name: "missing keyfile",
			setup: func() {
				gcsKeyfile = ""
				os.Unsetenv("GOOGLE_CLOUD_PROJECT")
			},
			teardown: func() {
				gcsKeyfile = ""
			},
			wantErr: true,
		},
		{
			name: "empty keyfile flag",
			setup: func() {
				gcsKeyfile = ""
				os.Setenv("GOOGLE_CLOUD_PROJECT", "")
			},
			teardown: func() {
				gcsKeyfile = ""
				os.Unsetenv("GOOGLE_CLOUD_PROJECT")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()

			err := validateGCSConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGCSConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGCSCommandSetup(t *testing.T) {
	// Test that GCS commands are properly initialized
	if gcsCmd == nil {
		t.Error("gcsCmd should be initialized")
	}

	if gcsPushCmd == nil {
		t.Error("gcsPushCmd should be initialized")
	}

	if gcsPullCmd == nil {
		t.Error("gcsPullCmd should be initialized")
	}

	// Test command properties
	if gcsCmd.Use != "gcs" {
		t.Errorf("gcsCmd.Use = %q, want %q", gcsCmd.Use, "gcs")
	}

	if gcsCmd.Short != "Google Cloud Storage operations" {
		t.Errorf("gcsCmd.Short = %q, want %q", gcsCmd.Short, "Google Cloud Storage operations")
	}

	if gcsPushCmd.Use != "push <bucket>/<image-path>:<tag>" {
		t.Errorf("gcsPushCmd.Use = %q, want %q", gcsPushCmd.Use, "push <bucket>/<image-path>:<tag>")
	}

	if gcsPullCmd.Use != "pull <bucket>/<image-path>:<tag>" {
		t.Errorf("gcsPullCmd.Use = %q, want %q", gcsPullCmd.Use, "pull <bucket>/<image-path>:<tag>")
	}
}
