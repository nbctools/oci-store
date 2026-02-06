package main

import (
	_ "github.com/distribution/distribution/v3/registry/storage/driver/gcs"
	"github.com/spf13/cobra"
)

var (
	gcsKeyfile       string
	gcsRootDirectory string
)

var gcsPushCmd = &cobra.Command{
	Use:   "push <bucket>/<image-path>:<tag>",
	Short: "Push a Docker image to Google Cloud Storage",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateGCSConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		localImage, _ := cmd.Flags().GetString("image")
		return pushImage(cmd.Context(), "gcs", args[0], localImage)
	},
}

var gcsPullCmd = &cobra.Command{
	Use:   "pull <bucket>/<image-path>:<tag>",
	Short: "Pull a Docker image from Google Cloud Storage",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateGCSConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return pullImage(cmd.Context(), "gcs", args[0])
	},
}

func validateGCSConfig() error {
	return nil
}

func init() {
	gcsCmd.AddCommand(gcsPushCmd, gcsPullCmd)

	gcsCmd.PersistentFlags().StringVar(&gcsKeyfile, "keyfile", "", "GCS keyfile")
	gcsCmd.PersistentFlags().StringVar(&gcsRootDirectory, "root-dir", "", "Root directory in GCS bucket (optional)")

	gcsPushCmd.Flags().StringP("image", "i", "", "Local Docker image to push (defaults to image-path:tag)")
}
