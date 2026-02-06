package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/distribution/distribution/v3/registry/storage/driver/s3-aws"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func pushImage(ctx context.Context, storageType string, storageRef string, localImage string) (err error) {
	backend, err := NewBackend(storageType)
	if err != nil {
		return err
	}

	ref, err := backend.ParseRef(storageRef)
	if err != nil {
		return err
	}

	if localImage == "" {
		localImage = storageRef
	}

	slog.Info("Pushing image", "image", localImage, "dest", fmt.Sprintf("%s://%s/%s:%s", ref.Type, ref.Bucket, ref.Path, ref.Tag), "bucket", ref.Bucket)
	regAddr, err := startRegistry(ctx, backend, ref.Bucket)
	if err != nil {
		return err
	}

	targetRef := fmt.Sprintf("%s/%s:%s", regAddr, ref.Path, ref.Tag)
	slog.Info("Target image reference", "ref", targetRef)
	slog.Info("Loading image from local Docker daemon", "source_image", localImage)

	localRef, err := name.ParseReference(localImage)
	if err != nil {
		return err
	}
	img, err := daemon.Image(localRef)
	if err != nil {
		return fmt.Errorf("failed to load image '%s' from local Docker daemon: %w", localImage, err)
	}
	slog.Info("Pushing image directly to target registry", "target", targetRef)
	dest, err := name.ParseReference(targetRef, name.Insecure) // Tell ParseReference that this might be insecure
	if err != nil {
		return fmt.Errorf("failed to parse target reference %s: %w", targetRef, err)
	}

	err = remote.Write(dest, img) // Push the image
	if err != nil {
		return fmt.Errorf("failed to push image directly to registry %s: %w", targetRef, err)
	}
	slog.Info("Image pushed directly to registry successfully!", "target", targetRef)
	return nil
}
