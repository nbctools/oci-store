package main

import (
	"fmt"
)

func NewBackend(storageType string) (StorageBackend, error) {
	switch storageType {
	case "s3":
		return newS3Backend(), nil
	case "gcs":
		return newGCSBackend(), nil
	case "azure":
		return newAzureBackend(), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

type S3Backend struct {
	RootDir   string
	Region    string
	Endpoint  string
	AccessKey string /* #nosec G117 */
	SecretKey string
}

func newS3Backend() *S3Backend {
	return &S3Backend{
		RootDir:   s3RootDirectory,
		Region:    s3Region,
		Endpoint:  s3Endpoint,
		AccessKey: s3AccessKey,
		SecretKey: s3SecretKey,
	}
}

func (s *S3Backend) Type() string {
	return "s3"
}

func (s *S3Backend) ParseRef(ref string) (*StorageRef, error) {
	return ParseStorageRef(ref, "s3")
}

func (s *S3Backend) GetStorageConfig(bucket string) map[string]interface{} {
	config := map[string]interface{}{
		"bucket": bucket,
	}

	if s.Region != "" {
		config["region"] = s.Region
	}
	if s.Endpoint != "" {
		config["regionendpoint"] = s.Endpoint
	}
	if s.AccessKey != "" {
		config["accesskey"] = s.AccessKey
	}
	if s.SecretKey != "" {
		config["secretkey"] = s.SecretKey
	}
	if s.RootDir != "" {
		config["rootdirectory"] = s.RootDir
	}

	loglevel := "error"
	if verbose {
		loglevel = "info"
	}
	config["loglevel"] = loglevel

	return config
}

func (s *S3Backend) ValidateConfig() error {
	if s.Region == "" {
		return fmt.Errorf("S3 requires region to be specified")
	}
	return nil
}

type GCSBackend struct {
	RootDir string
	Keyfile string // For GCS
}

func newGCSBackend() *GCSBackend {
	return &GCSBackend{
		RootDir: gcsRootDirectory,
		Keyfile: gcsKeyfile,
	}
}

func (g *GCSBackend) Type() string {
	return "gcs"
}

func (g *GCSBackend) ParseRef(ref string) (*StorageRef, error) {
	return ParseStorageRef(ref, "gcs")
}

func (g *GCSBackend) GetStorageConfig(bucket string) map[string]interface{} {
	config := map[string]interface{}{
		"bucket": bucket,
	}

	if g.Keyfile != "" {
		config["keyfile"] = g.Keyfile
	}
	if g.RootDir != "" {
		config["rootdirectory"] = g.RootDir
	}

	loglevel := "error"
	if verbose {
		loglevel = "info"
	}
	config["loglevel"] = loglevel

	return config
}

func (g *GCSBackend) ValidateConfig() error {
	return nil
}

type AzureBackend struct {
	AccountName    string // For Azure
	AccountKey     string // For Azure
	Container      string // For Azure
	CredentialType string //For Azure
	RootDir        string
	Secret         string /* #nosec G117 */
	TenantID       string
	ClientID       string
}

func newAzureBackend() *AzureBackend {
	return &AzureBackend{
		AccountName:    azureAccountName,
		AccountKey:     azureAccountKey,
		CredentialType: azureCredentialType,
		RootDir:        azureRootDirectory,
		Secret:         azureSecret,
		ClientID:       azureClientId,
		TenantID:       azureTenantId,
	}
}

func (a *AzureBackend) Type() string {
	return "azure"
}

func (a *AzureBackend) ParseRef(ref string) (*StorageRef, error) {
	return ParseStorageRef(ref, "azure")
}

func (a *AzureBackend) GetStorageConfig(container string) map[string]interface{} {
	config := map[string]interface{}{
		"container": container,
	}

	if a.AccountName != "" {
		config["accountname"] = a.AccountName
	}
	if a.AccountKey != "" {
		config["accountkey"] = a.AccountKey
	}
	if a.RootDir != "" {
		config["rootdirectory"] = a.RootDir
	}
	if a.CredentialType != "" {
		config["credentials"] = map[string]string{"type": a.CredentialType,
			"secret":   a.Secret,
			"clientid": a.ClientID,
			"tenantid": a.TenantID}
	}

	loglevel := "error"
	if verbose {
		loglevel = "info"
	}
	config["loglevel"] = loglevel

	return config
}

func (a *AzureBackend) ValidateConfig() error {
	if a.AccountName == "" {
		return fmt.Errorf("missing required Azure account name")
	}
	if a.AccountKey == "" {
		return fmt.Errorf("missing required account key")
	}
	return nil
}
