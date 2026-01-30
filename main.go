package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	region        string
	endpoint      string
	accessKey     string
	secretKey     string
	rootDirectory string
	verbose       bool
	logopts       = &slog.HandlerOptions{}
)

var rootCmd = &cobra.Command{
	Use:   "oci-s3",
	Short: "Push and pull Docker images to/from S3",
	Long:  `A CLI tool to store Docker/OCI images in S3 using distribution registry's S3 storage driver`,
}

var pushCmd = &cobra.Command{
	Use:   "push <s3-bucket>/<image-path>:<tag>",
	Short: "Push a Docker image to S3",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return valRegion()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		localImage, _ := cmd.Flags().GetString("image")
		return pushImage(cmd.Context(), args[0], localImage)
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull <s3-bucket>/<image-path>:<tag>",
	Short: "Pull a Docker image from S3",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return valRegion()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return pullImage(cmd.Context(), args[0])
	},
}

func valRegion() error {
	if region == "" {
		if region = getEnv("AWS_REGION"); region == "" {
			return errors.New("AWS region not specified")
		}
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region (defaults to AWS_REGION env var)")
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "S3-compatible endpoint (optional)")
	rootCmd.PersistentFlags().StringVar(&accessKey, "access-key", "", "AWS access key (defaults to AWS_ACCESS_KEY_ID env var)")
	rootCmd.PersistentFlags().StringVar(&secretKey, "secret-key", "", "AWS secret key (defaults to AWS_SECRET_ACCESS_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&rootDirectory, "root-dir", "", "Root directory in S3 bucket (optional)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	pushCmd.Flags().StringP("image", "i", "", "Local Docker image to push (defaults to image-path:tag)")

	rootCmd.AddCommand(pushCmd, pullCmd)
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
