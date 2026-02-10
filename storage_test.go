package main

import (
	"testing"
)

func TestParseStorageRef(t *testing.T) {
	tests := []struct {
		name        string
		ref         string
		storageType string
		want        *StorageRef
		wantErr     bool
	}{
		{
			name:        "valid S3 ref",
			ref:         "my-bucket/path/to/image:v1.0",
			storageType: "s3",
			want: &StorageRef{
				Bucket: "my-bucket",
				Path:   "path/to/image",
				Tag:    "v1.0",
				Type:   "s3",
			},
			wantErr: false,
		},
		{
			name:        "valid GCS ref",
			ref:         "my-gcs-bucket/app:v2.0",
			storageType: "gcs",
			want: &StorageRef{
				Bucket: "my-gcs-bucket",
				Path:   "app",
				Tag:    "v2.0",
				Type:   "gcs",
			},
			wantErr: false,
		},
		{
			name:        "valid Azure ref",
			ref:         "my-container/service:v3.0",
			storageType: "azure",
			want: &StorageRef{
				Bucket: "my-container",
				Path:   "service",
				Tag:    "v3.0",
				Type:   "azure",
			},
			wantErr: false,
		},
		{
			name:        "invalid ref - no slash",
			ref:         "invalid-ref",
			storageType: "s3",
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "invalid ref - no tag",
			ref:         "bucket/path",
			storageType: "s3",
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStorageRef(tt.ref, tt.storageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStorageRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("ParseStorageRef() returned nil, want %v", tt.want)
					return
				}
				if got.Bucket != tt.want.Bucket || got.Path != tt.want.Path || got.Tag != tt.want.Tag || got.Type != tt.want.Type {
					t.Errorf("ParseStorageRef() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestNewBackend(t *testing.T) {
	tests := []struct {
		name        string
		storageType string
		wantErr     bool
	}{
		{
			name:        "valid S3",
			storageType: "s3",
			wantErr:     false,
		},
		{
			name:        "valid GCS",
			storageType: "gcs",
			wantErr:     false,
		},
		{
			name:        "valid Azure",
			storageType: "azure",
			wantErr:     false,
		},
		{
			name:        "invalid storage type",
			storageType: "invalid",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBackend(tt.storageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestS3Backend(t *testing.T) {
	// Set required config for validation
	oldRegion := s3Region
	s3Region = "us-east-1"
	defer func() { s3Region = oldRegion }()

	backend, err := NewBackend("s3")
	if err != nil {
		t.Fatalf("NewBackend() error = %v", err)
	}

	s3Backend, ok := backend.(*S3Backend)
	if !ok {
		t.Fatalf("Expected S3Backend, got %T", backend)
	}

	// Test Type method
	if s3Backend.Type() != "s3" {
		t.Errorf("S3Backend.Type() = %q, want %q", s3Backend.Type(), "s3")
	}

	// Test ParseRef
	ref, err := s3Backend.ParseRef("my-bucket/app:v1.0")
	if err != nil {
		t.Errorf("S3Backend.ParseRef() error = %v", err)
	}
	if ref.Type != "s3" || ref.Bucket != "my-bucket" || ref.Path != "app" || ref.Tag != "v1.0" {
		t.Errorf("S3Backend.ParseRef() = %+v, want Type=s3, Bucket=my-bucket, Path=app, Tag=v1.0", ref)
	}

	// Test GetStorageConfig
	storageConfig := s3Backend.GetStorageConfig("test-bucket")
	if storageConfig["bucket"] != "test-bucket" {
		t.Errorf("S3Backend.GetStorageConfig() = %+v, want bucket=test-bucket", storageConfig)
	}
	if storageConfig["region"] != "us-east-1" {
		t.Errorf("S3Backend.GetStorageConfig() region = %q, want us-east-1", storageConfig["region"])
	}

	// Test ValidateConfig
	if err := s3Backend.ValidateConfig(); err != nil {
		t.Errorf("S3Backend.ValidateConfig() error = %v", err)
	}

	// Test validation failure - create new backend with empty region
	oldRegion = s3Region
	s3Region = ""
	invalidBackend, _ := NewBackend("s3")
	invalidS3Backend := invalidBackend.(*S3Backend)
	if err := invalidS3Backend.ValidateConfig(); err == nil {
		t.Error("S3Backend.ValidateConfig() should have failed with empty region")
	}
	s3Region = oldRegion
}

func TestGCSBackend(t *testing.T) {
	backend, err := NewBackend("gcs")
	if err != nil {
		t.Fatalf("NewBackend() error = %v", err)
	}

	gcsBackend, ok := backend.(*GCSBackend)
	if !ok {
		t.Fatalf("Expected GCSBackend, got %T", backend)
	}

	// Test Type method
	if gcsBackend.Type() != "gcs" {
		t.Errorf("GCSBackend.Type() = %q, want %q", gcsBackend.Type(), "gcs")
	}

	// Test ParseRef
	ref, err := gcsBackend.ParseRef("my-bucket/app:v1.0")
	if err != nil {
		t.Errorf("GCSBackend.ParseRef() error = %v", err)
	}
	if ref.Type != "gcs" || ref.Bucket != "my-bucket" || ref.Path != "app" || ref.Tag != "v1.0" {
		t.Errorf("GCSBackend.ParseRef() = %+v, want Type=gcs, Bucket=my-bucket, Path=app, Tag=v1.0", ref)
	}

	// Test GetStorageConfig
	storageConfig := gcsBackend.GetStorageConfig("test-bucket")
	if storageConfig["bucket"] != "test-bucket" {
		t.Errorf("GCSBackend.GetStorageConfig() = %+v, want bucket=test-bucket", storageConfig)
	}

	// Test with keyfile set
	oldKeyfile := gcsKeyfile
	gcsKeyfile = "test-project.json"
	defer func() { gcsKeyfile = oldKeyfile }()

	// Recreate backend to pick up new keyfile
	backend, err = NewBackend("gcs")
	if err != nil {
		t.Fatalf("NewBackend() error = %v", err)
	}
	gcsBackend = backend.(*GCSBackend)

	storageConfig = gcsBackend.GetStorageConfig("test-bucket")
	if storageConfig["keyfile"] != "test-project.json" {
		t.Errorf("GCSBackend.GetStorageConfig() keyfile = %q, want test-project.json", storageConfig["keyfile"])
	}
}

func TestAzureBackend(t *testing.T) {
	// Set required config for validation
	oldAccountName := azureAccountName
	oldAccountKey := azureAccountKey
	azureAccountName = "test-account"
	azureAccountKey = "test-key"
	defer func() {
		azureAccountName = oldAccountName
		azureAccountKey = oldAccountKey
	}()

	backend, err := NewBackend("azure")
	if err != nil {
		t.Fatalf("NewBackend() error = %v", err)
	}

	azureBackend, ok := backend.(*AzureBackend)
	if !ok {
		t.Fatalf("Expected AzureBackend, got %T", backend)
	}

	// Test Type method
	if azureBackend.Type() != "azure" {
		t.Errorf("AzureBackend.Type() = %q, want %q", azureBackend.Type(), "azure")
	}

	// Test ParseRef
	ref, err := azureBackend.ParseRef("my-container/app:v1.0")
	if err != nil {
		t.Errorf("AzureBackend.ParseRef() error = %v", err)
	}
	if ref.Type != "azure" || ref.Bucket != "my-container" || ref.Path != "app" || ref.Tag != "v1.0" {
		t.Errorf("AzureBackend.ParseRef() = %+v, want Type=azure, Bucket=my-container, Path=app, Tag=v1.0", ref)
	}

	// Test GetStorageConfig
	storageConfig := azureBackend.GetStorageConfig("test-container")
	if storageConfig["container"] != "test-container" {
		t.Errorf("AzureBackend.GetStorageConfig() = %+v, want container=test-container", storageConfig)
	}
	if storageConfig["accountname"] != "test-account" {
		t.Errorf("AzureBackend.GetStorageConfig() accountname = %q, want test-account", storageConfig["accountname"])
	}

	// Test ValidateConfig
	if err := azureBackend.ValidateConfig(); err != nil {
		t.Errorf("AzureBackend.ValidateConfig() error = %v", err)
	}

	// Test validation failure - missing account name
	oldAccountName = azureAccountName
	azureAccountName = ""
	invalidBackend, _ := NewBackend("azure")
	invalidAzureBackend := invalidBackend.(*AzureBackend)
	if err := invalidAzureBackend.ValidateConfig(); err == nil {
		t.Error("AzureBackend.ValidateConfig() should have failed with empty account name")
	}
	azureAccountName = oldAccountName

	// Test validation failure - missing account key
	oldAccountKey = azureAccountKey
	azureAccountKey = ""
	invalidBackend, _ = NewBackend("azure")
	invalidAzureBackend = invalidBackend.(*AzureBackend)
	if err := invalidAzureBackend.ValidateConfig(); err == nil {
		t.Error("AzureBackend.ValidateConfig() should have failed with empty account key")
	}
	azureAccountKey = oldAccountKey
}
