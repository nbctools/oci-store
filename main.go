package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	logopts = &slog.HandlerOptions{}
)

var rootCmd = &cobra.Command{
	Use:   "oci-store",
	Short: "Push and pull Docker images directly to/from cloud storage",
	Long:  `A CLI tool to store Docker/OCI images in cloud storage using distribution registry's storage drivers`,
}

var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "S3 storage operations",
}

var gcsCmd = &cobra.Command{
	Use:   "gcs",
	Short: "Google Cloud Storage operations",
}

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Azure Blob Storage operations",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	rootCmd.AddCommand(s3Cmd, gcsCmd, azureCmd)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			logopts.Level = slog.LevelDebug
		}
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, logopts)))
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		slog.Error("Command failed", "error", err)
		os.Exit(1)
	}
}
