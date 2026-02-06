package main

import (
	"fmt"
	"strings"
)

type StorageRef struct {
	Bucket string
	Path   string
	Tag    string
	Type   string
}

type StorageBackend interface {
	Type() string
	ParseRef(ref string) (*StorageRef, error)
	GetStorageConfig(bucket string) map[string]interface{}
	ValidateConfig() error
}

func ParseStorageRef(ref string, storageType string) (*StorageRef, error) {
	bucket, pathTag, ok := strings.Cut(ref, "/")
	if !ok {
		return nil, fmt.Errorf("invalid %s reference format, expected: bucket/path:tag", storageType)
	}
	path, tag, ok := strings.Cut(pathTag, ":")
	if !ok {
		return nil, fmt.Errorf("missing tag in reference")
	}

	return &StorageRef{
		Bucket: bucket,
		Path:   path,
		Tag:    tag,
		Type:   storageType,
	}, nil
}
