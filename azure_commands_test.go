package main

import (
	"os"
	"testing"
)

func TestValidateAzureConfig(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		wantErr  bool
	}{
		{
			name: "valid config with flags",
			setup: func() {
				azureAccountName = "testaccount"
				azureAccountKey = "testkey"
			},
			teardown: func() {
				azureAccountName = ""
				azureAccountKey = ""
			},
			wantErr: false,
		},
		{
			name: "valid config with env vars",
			setup: func() {
				azureAccountName = ""
				azureAccountKey = ""
				os.Setenv("AZURE_STORAGE_ACCOUNT", "envaccount")
				os.Setenv("AZURE_STORAGE_KEY", "envkey")
			},
			teardown: func() {
				azureAccountName = ""
				azureAccountKey = ""
				os.Unsetenv("AZURE_STORAGE_ACCOUNT")
				os.Unsetenv("AZURE_STORAGE_KEY")
			},
			wantErr: false,
		},
		{
			name: "missing account name",
			setup: func() {
				azureAccountName = ""
				azureAccountKey = "testkey"
				os.Unsetenv("AZURE_STORAGE_ACCOUNT")
				os.Unsetenv("AZURE_STORAGE_KEY")
			},
			teardown: func() {
				azureAccountName = ""
				azureAccountKey = ""
			},
			wantErr: true,
		},
		{
			name: "missing account key",
			setup: func() {
				azureAccountName = "testaccount"
				azureAccountKey = ""
				os.Unsetenv("AZURE_STORAGE_ACCOUNT")
				os.Unsetenv("AZURE_STORAGE_KEY")
			},
			teardown: func() {
				azureAccountName = ""
				azureAccountKey = ""
			},
			wantErr: true,
		},
		{
			name: "missing both",
			setup: func() {
				azureAccountName = ""
				azureAccountKey = ""
				os.Unsetenv("AZURE_STORAGE_ACCOUNT")
				os.Unsetenv("AZURE_STORAGE_KEY")
			},
			teardown: func() {
				azureAccountName = ""
				azureAccountKey = ""
			},
			wantErr: true,
		},
		{
			name: "empty account name with env key",
			setup: func() {
				azureAccountName = ""
				azureAccountKey = ""
				os.Setenv("AZURE_STORAGE_ACCOUNT", "")
				os.Setenv("AZURE_STORAGE_KEY", "testkey")
			},
			teardown: func() {
				azureAccountName = ""
				azureAccountKey = ""
				os.Unsetenv("AZURE_STORAGE_ACCOUNT")
				os.Unsetenv("AZURE_STORAGE_KEY")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()

			err := validateAzureConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAzureConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAzureCommandSetup(t *testing.T) {
	// Test that Azure commands are properly initialized
	if azureCmd == nil {
		t.Error("azureCmd should be initialized")
	}

	if azurePushCmd == nil {
		t.Error("azurePushCmd should be initialized")
	}

	if azurePullCmd == nil {
		t.Error("azurePullCmd should be initialized")
	}

	// Test command properties
	if azureCmd.Use != "azure" {
		t.Errorf("azureCmd.Use = %q, want %q", azureCmd.Use, "azure")
	}

	if azureCmd.Short != "Azure Blob Storage operations" {
		t.Errorf("azureCmd.Short = %q, want %q", azureCmd.Short, "Azure Blob Storage operations")
	}

	if azurePushCmd.Use != "push <container>/<image-path>:<tag>" {
		t.Errorf("azurePushCmd.Use = %q, want %q", azurePushCmd.Use, "push <container>/<image-path>:<tag>")
	}

	if azurePullCmd.Use != "pull <container>/<image-path>:<tag>" {
		t.Errorf("azurePullCmd.Use = %q, want %q", azurePullCmd.Use, "pull <container>/<image-path>:<tag>")
	}
}
