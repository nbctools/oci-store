package main

import (
	"os"
	"testing"
)

func TestValidateS3Config(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		wantErr  bool
	}{
		{
			name: "valid config with region flag",
			setup: func() {
				s3Region = "us-east-1"
			},
			teardown: func() {
				s3Region = ""
			},
			wantErr: false,
		},
		{
			name: "valid config with env var",
			setup: func() {
				s3Region = ""
				os.Setenv("AWS_REGION", "us-west-2")
			},
			teardown: func() {
				s3Region = ""
				os.Unsetenv("AWS_REGION")
			},
			wantErr: false,
		},
		{
			name: "missing region",
			setup: func() {
				s3Region = ""
				os.Unsetenv("AWS_REGION")
			},
			teardown: func() {
				s3Region = ""
			},
			wantErr: true,
		},
		{
			name: "empty region flag",
			setup: func() {
				s3Region = ""
				os.Setenv("AWS_REGION", "")
			},
			teardown: func() {
				s3Region = ""
				os.Unsetenv("AWS_REGION")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()

			err := validateS3Config()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateS3Config() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestS3CommandSetup(t *testing.T) {
	// Test that S3 commands are properly initialized
	if s3Cmd == nil {
		t.Error("s3Cmd should be initialized")
	}

	if s3PushCmd == nil {
		t.Error("s3PushCmd should be initialized")
	}

	if s3PullCmd == nil {
		t.Error("s3PullCmd should be initialized")
	}

	// Test command properties
	if s3Cmd.Use != "s3" {
		t.Errorf("s3Cmd.Use = %q, want %q", s3Cmd.Use, "s3")
	}

	if s3Cmd.Short != "S3 storage operations" {
		t.Errorf("s3Cmd.Short = %q, want %q", s3Cmd.Short, "S3 storage operations")
	}

	if s3PushCmd.Use != "push <bucket>/<image-path>:<tag>" {
		t.Errorf("s3PushCmd.Use = %q, want %q", s3PushCmd.Use, "push <bucket>/<image-path>:<tag>")
	}

	if s3PullCmd.Use != "pull <bucket>/<image-path>:<tag>" {
		t.Errorf("s3PullCmd.Use = %q, want %q", s3PullCmd.Use, "pull <bucket>/<image-path>:<tag>")
	}
}
