package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
)

func pullImage(ctx context.Context, s3Ref string) error {
	ref, err := parseS3Ref(s3Ref)
	if err != nil {
		return err
	}
	slog.Info("Pulling image", "bucket", ref.Bucket, "image", ref.Path+":"+ref.Tag)

	regAddr, err := startRegistry(ctx, ref.Bucket)
	if err != nil {
		return err
	}
	srcRef := fmt.Sprintf("%s/%s:%s", regAddr, ref.Path, ref.Tag)
	img, err := crane.Pull(srcRef, crane.Insecure)
	if err != nil {
		return err
	}
	tag, err := name.NewTag(s3Ref)
	if err != nil {
		return err
	}
	_, err = daemon.Write(tag, img)
	slog.Info("Image pulled", "name", s3Ref)
	return err
}
