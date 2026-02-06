package main

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

var (
	s3Region        string
	s3Endpoint      string
	s3AccessKey     string
	s3SecretKey     string
	s3RootDirectory string
)

var s3PushCmd = &cobra.Command{
	Use:   "push <bucket>/<image-path>:<tag>",
	Short: "Push a Docker image to S3",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateS3Config()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		localImage, _ := cmd.Flags().GetString("image")
		return pushImage(cmd.Context(), "s3", args[0], localImage)
	},
}

var s3PullCmd = &cobra.Command{
	Use:   "pull <bucket>/<image-path>:<tag>",
	Short: "Pull a Docker image from S3",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateS3Config()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return pullImage(cmd.Context(), "s3", args[0])
	},
}

func validateS3Config() error {
	if s3Region == "" {
		if s3Region = strings.TrimSpace(getEnv("AWS_REGION")); s3Region == "" {
			return errors.New("S3 requires region to be specified via --region or AWS_REGION env var")
		}
	}
	if s3AccessKey == "" {
		s3AccessKey = strings.TrimSpace(getEnv("AWS_ACCESS_KEY_ID"))
	}
	if s3SecretKey == "" {
		s3SecretKey = strings.TrimSpace(getEnv("AWS_SECRET_ACCESS_KEY"))
	}
	return nil
}

func init() {
	s3Cmd.AddCommand(s3PushCmd, s3PullCmd)

	s3Cmd.PersistentFlags().StringVarP(&s3Region, "region", "r", "", "AWS region (defaults to AWS_REGION env var)")
	s3Cmd.PersistentFlags().StringVarP(&s3Endpoint, "endpoint", "e", "", "S3-compatible endpoint (optional)")
	s3Cmd.PersistentFlags().StringVar(&s3AccessKey, "access-key", "", "AWS access key (defaults to AWS_ACCESS_KEY_ID env var)")
	s3Cmd.PersistentFlags().StringVar(&s3SecretKey, "secret-key", "", "AWS secret key (defaults to AWS_SECRET_ACCESS_KEY env var)")
	s3Cmd.PersistentFlags().StringVar(&s3RootDirectory, "root-dir", "", "Root directory in S3 bucket (optional)")

	s3PushCmd.Flags().StringP("image", "i", "", "Local Docker image to push (defaults to image-path:tag)")
}
