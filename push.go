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

type S3Ref struct {
	Bucket string
	Path   string
	Tag    string
}

func parseS3Ref(ref string) (*S3Ref, error) {
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid S3 reference format, expected: bucket/path:tag")
	}

	bucket := parts[0]
	pathTag := parts[1]

	tagParts := strings.SplitN(pathTag, ":", 2)
	if len(tagParts) != 2 {
		return nil, fmt.Errorf("missing tag in reference")
	}

	return &S3Ref{
		Bucket: bucket,
		Path:   tagParts[0],
		Tag:    tagParts[1],
	}, nil
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func pushImage(ctx context.Context, s3Ref string, localImage string) (err error) {
	ref, err := parseS3Ref(s3Ref)
	if err != nil {
		return err
	}

	if localImage == "" {
		localImage = s3Ref
	}

	slog.Info("Pushing image", "image", localImage, "dest", fmt.Sprintf("s3://%s/%s:%s", ref.Bucket, ref.Path, ref.Tag), "bucket", ref.Bucket)
	regAddr, err := startRegistry(ctx, ref.Bucket)
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
