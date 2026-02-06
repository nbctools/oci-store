package main

import (
	"errors"

	_ "github.com/distribution/distribution/v3/registry/storage/driver/azure"
	"github.com/spf13/cobra"
)

var (
	azureAccountName    string
	azureAccountKey     string
	azureRootDirectory  string
	azureCredentialType string
	azureSecret         string
	azureTenantId       string
	azureClientId       string
)

var azurePushCmd = &cobra.Command{
	Use:   "push <container>/<image-path>:<tag>",
	Short: "Push a Docker image to Azure Blob Storage",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateAzureConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		localImage, _ := cmd.Flags().GetString("image")
		return pushImage(cmd.Context(), "azure", args[0], localImage)
	},
}

var azurePullCmd = &cobra.Command{
	Use:   "pull <container>/<image-path>:<tag>",
	Short: "Pull a Docker image from Azure Blob Storage",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateAzureConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return pullImage(cmd.Context(), "azure", args[0])
	},
}

func validateAzureConfig() error {
	if azureAccountName == "" {
		if azureAccountName = getEnv("AZURE_STORAGE_ACCOUNT"); azureAccountName == "" {
			return errors.New("account name needs to be specified via --account-name or AZURE_STORAGE_ACCOUNT env var")
		}
	}
	if azureAccountKey == "" {
		if azureAccountKey = getEnv("AZURE_STORAGE_KEY"); azureAccountKey == "" {
			return errors.New("account key needs to be specified via --account-key or AZURE_STORAGE_KEY env var")
		}
	}
	if azureClientId == "" {
		azureClientId = getEnv("AZURE_CLIENT_ID")
	}
	if azureTenantId == "" {
		azureTenantId = getEnv("AZURE_TENANT_ID")
	}
	if azureSecret == "" {
		azureSecret = getEnv("AZURE_SECRET")
	}
	return nil
}

func init() {
	azureCmd.AddCommand(azurePushCmd, azurePullCmd)

	azureCmd.PersistentFlags().StringVarP(&azureAccountName, "account-name", "a", "", "Azure storage account name (defaults to AZURE_STORAGE_ACCOUNT env var)")
	azureCmd.PersistentFlags().StringVarP(&azureAccountKey, "account-key", "k", "", "Azure storage account key (defaults to AZURE_STORAGE_KEY env var)")
	azureCmd.PersistentFlags().StringVar(&azureRootDirectory, "root-dir", "", "Root directory in Azure container (optional)")
	azureCmd.PersistentFlags().StringVar(&azureCredentialType, "credential-type", "client_secret", "Azure credentials used to authenticate with Azure blob storage")
	azureCmd.PersistentFlags().StringVar(&azureClientId, "client-id", "", "The unique app ID (defaults to AZURE_CLIENT_ID)")
	azureCmd.PersistentFlags().StringVar(&azureTenantId, "tenant-id", "", "The directory(tenant) ID (defaults to AZURE_TENANT_ID)")
	azureCmd.PersistentFlags().StringVar(&azureSecret, "secret", "", "The client secret(defaults to AZURE_SECRET)")

	azurePushCmd.Flags().StringP("image", "i", "", "Local Docker image to push (defaults to image-path:tag)")
}
