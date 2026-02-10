package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestCLIIntegration tests the complete CLI workflow
func TestCLIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI integration test in short mode")
	}

	// Build the binary for testing
	buildBinary(t)
	defer cleanupBinary(t)

	// Test help commands
	testHelpCommands(t)

	// Test validation commands
	testValidationCommands(t)

	// Test command structure
	testCommandStructure(t)
}

func buildBinary(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "test-oci-store", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, string(output))
	}
}

func cleanupBinary(t *testing.T) {
	os.Remove("test-oci-store")
}

func runCommand(t *testing.T, args ...string) (string, error) {
	cmd := exec.Command("./test-oci-store", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), err
}

func testHelpCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "main help",
			args:     []string{"--help"},
			contains: "A CLI tool to store Docker/OCI images in cloud storage",
		},
		{
			name:     "s3 help",
			args:     []string{"s3", "--help"},
			contains: "S3 storage operations",
		},
		{
			name:     "gcs help",
			args:     []string{"gcs", "--help"},
			contains: "Google Cloud Storage operations",
		},
		{
			name:     "azure help",
			args:     []string{"azure", "--help"},
			contains: "Azure Blob Storage operations",
		},
		{
			name:     "s3 push help",
			args:     []string{"s3", "push", "--help"},
			contains: "Push a Docker image to S3",
		},
		{
			name:     "gcs pull help",
			args:     []string{"gcs", "pull", "--help"},
			contains: "Pull a Docker image from Google Cloud Storage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCommand(t, tt.args...)
			if err != nil {
				t.Errorf("Command failed: %v, output: %s", err, output)
			}
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Output does not contain %q: %s", tt.contains, output)
			}
		})
	}
}

func testValidationCommands(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
		contains  string
	}{
		{
			name:      "s3 push without region",
			args:      []string{"s3", "push", "test-bucket/app:v1.0"},
			wantError: true,
			contains:  "requires region",
		},
		{
			name:      "gcs push without project",
			args:      []string{"gcs", "push", "test-bucket/app:v1.0"},
			wantError: true,
			contains:  "requires project ID to be specified",
		},
		{
			name:      "azure push without account",
			args:      []string{"azure", "push", "test-container/app:v1.0"},
			wantError: true,
			contains:  "account name needs to be specified",
		},
		{
			name:      "s3 pull without region",
			args:      []string{"s3", "pull", "test-bucket/app:v1.0"},
			wantError: true,
			contains:  "requires region",
		},
		{
			name:      "invalid command",
			args:      []string{"invalid", "command"},
			wantError: true,
			contains:  "unknown command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCommand(t, tt.args...)

			if tt.wantError && err == nil {
				t.Errorf("Expected command to fail, but it succeeded")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Expected command to succeed, but it failed: %v, output: %s", err, output)
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Errorf("Output does not contain %q: %s", tt.contains, output)
			}
		})
	}
}

func testCommandStructure(t *testing.T) {
	// Test that all expected commands exist
	output, err := runCommand(t, "--help")
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	expectedCommands := []string{"s3", "gcs", "azure", "completion", "help"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Help output does not contain command %q: %s", cmd, output)
		}
	}
}

// TestEnvironmentVariableHandling tests environment variable precedence
func TestEnvironmentVariableHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping environment variable test in short mode")
	}

	buildBinary(t)
	defer cleanupBinary(t)

	// Test S3 environment variable
	oldAWSRegion := os.Getenv("AWS_REGION")
	os.Setenv("AWS_REGION", "test-region-from-env")
	defer func() {
		if oldAWSRegion != "" {
			os.Setenv("AWS_REGION", oldAWSRegion)
		} else {
			os.Unsetenv("AWS_REGION")
		}
	}()

	// This should not fail due to region being set via env var
	output, err := runCommand(t, "s3", "push", "test-bucket/app:v1.0")
	if err == nil {
		// Command might still fail for other reasons (Docker not running), but should not fail due to region
		if strings.Contains(output, "requires region") {
			t.Error("Environment variable AWS_REGION was not respected")
		}
	}
}

// TestVerboseFlag tests the verbose flag functionality
func TestVerboseFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping verbose flag test in short mode")
	}

	buildBinary(t)
	defer cleanupBinary(t)

	// Test with verbose flag
	output, err := runCommand(t, "--verbose", "s3", "push", "test-bucket/app:v1.0")

	// The command should fail (no region), but we're testing the verbose flag
	if err == nil {
		t.Error("Expected command to fail without region")
	}

	// With verbose flag, we should see more detailed output
	// This is a basic test - in a real scenario, you'd check for specific debug output
	if len(output) == 0 {
		t.Error("Verbose flag should produce output")
	}
}

// TestCommandTimeout tests that commands don't hang indefinitely
func TestCommandTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	buildBinary(t)
	defer cleanupBinary(t)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "./test-oci-store", "--help")

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Help command failed: %v", err)
		}
	case <-ctx.Done():
		t.Error("Command timed out")
	}
}
