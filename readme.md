# OCI-S3: Docker Image Storage on S3

A CLI tool to push and pull Docker/OCI images to/from Amazon S3 with content-addressable storage optimization.

## Features

- **Distribution Registry Storage**: Uses the battle-tested storage driver from distribution/distribution
- **S3 Storage Driver**: Leverages the official S3 storage driver with all its optimizations
- **Content-Addressable Storage**: Automatic deduplication of layers across images
- **Production-Ready**: Built on the same storage backend used by Docker Registry
- **Multipart Uploads**: Efficient handling of large layers with S3 multipart uploads
- **Simple CLI**: Familiar push/pull interface

## Installation

```bash
go install github.com/yourusername/oci-s3@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/oci-s3
cd oci-s3
go build -o oci-s3
```

## Prerequisites

- Go 1.21 or later
- Docker installed and running
- AWS credentials configured (via `~/.aws/credentials`, environment variables, or IAM role)
- S3 bucket with appropriate permissions

The tool uses the same S3 storage driver as the Docker Registry, which means it supports:
- S3-compatible storage (MinIO, DigitalOcean Spaces, etc.)
- IAM roles for EC2 instances
- S3 encryption
- CloudFront integration (when configured)

## Usage

### Push an image to S3

```bash
# Push a local Docker image to S3
oci-s3 push my-bucket/myapp:v1.0

# Push with explicit local image name
oci-s3 push my-bucket/myapp:v1.0 --image localhost/myapp:latest

# Push to a specific region
oci-s3 push my-bucket/myapp:v1.0 --region us-west-2

# Push to S3-compatible storage (e.g., MinIO)
oci-s3 push my-bucket/myapp:v1.0 --endpoint http://localhost:9000

# Use a root directory in the bucket
oci-s3 push my-bucket/myapp:v1.0 --root-dir /registry
```

### Pull an image from S3

```bash
# Pull an image from S3 and load into Docker
oci-s3 pull my-bucket/myapp:v1.0
```

## S3 Storage Layout

The tool uses the same layout as Docker Distribution Registry:

```
s3://my-bucket/
  docker/
    registry/
      v2/
        repositories/
          myapp/
            _manifests/
              tags/
                v1.0/
                  current/link  → points to manifest digest
            _layers/
              sha256/
                abc123.../link  → points to blob
        blobs/
          sha256/
            ab/
              abc123.../data    (actual layer data)
            de/
              def456.../data
```

This layout provides:
- **Content-addressable storage**: Blobs stored by digest for deduplication
- **Repository isolation**: Each image has its own namespace
- **Atomic operations**: Tags are symbolic links for atomic updates
- **Garbage collection support**: Unused blobs can be identified and removed

## How It Works

This tool leverages the storage architecture from `distribution/distribution` (the official Docker Registry implementation):

1. **Storage Driver**: Uses the production-grade S3 storage driver that powers Docker Registry
2. **Registry Storage**: Implements the full registry storage API for proper blob and manifest management
3. **Content Addressable**: All blobs are stored by their SHA256 digest
4. **Optimizations**:
   - Multipart uploads for large layers (configurable chunk size)
   - Parallel multipart copy operations
   - Blob deduplication across all images
   - S3 transfer acceleration support
   - CloudFront CDN integration capability

**Push Process**:
1. Retrieves image from local Docker daemon
2. Uses registry storage API to upload blobs (with automatic deduplication)
3. Uploads manifest and creates tag reference
4. S3 driver handles multipart uploads for large blobs

**Pull Process**:
1. Resolves tag to manifest digest via registry storage
2. Downloads manifest and parses layer references
3. Downloads all blobs using S3 driver (with verification)
4. Reconstructs image and loads into Docker daemon

## AWS Permissions

Your AWS credentials need the following S3 permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::my-bucket/*",
        "arn:aws:s3:::my-bucket"
      ]
    }
  ]
}
```

## Benefits of Using Distribution Registry Storage

- **Battle-tested**: Same storage backend used by Docker Registry in production
- **Deduplication**: Identical layers shared across images are stored only once
- **S3 Optimizations**: 
  - Multipart uploads for large files
  - Configurable chunk sizes
  - S3 transfer acceleration
  - Reduced redundancy storage options
- **Integrity**: Content hash ensures data integrity at every layer
- **Atomic Updates**: Tag updates are atomic operations
- **Efficient Transfer**: Only missing layers need to be uploaded/downloaded
- **S3-Compatible**: Works with MinIO, DigitalOcean Spaces, and other S3-compatible storage

## Example Workflow

```bash
# Build an image
docker build -t myapp:v1.0 .

# Push to S3
oci-s3 push my-bucket/myapp:v1.0

# On another machine with same AWS credentials
oci-s3 pull my-bucket/myapp:v1.0

# Run the image
docker run myapp:v1.0
```

## Advanced Configuration

The tool supports all S3 storage driver options from Docker Registry:

### Environment Variables

```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
```

### Command Line Flags

```bash
# S3-compatible storage (MinIO, etc.)
oci-s3 push my-bucket/myapp:v1.0 \
  --endpoint http://localhost:9000 \
  --region us-east-1

# Use a subdirectory in the bucket
oci-s3 push my-bucket/myapp:v1.0 \
  --root-dir /docker/registry

# Combine options
oci-s3 pull my-bucket/myapp:v1.0 \
  --region eu-west-1 \
  --root-dir /production
```

### Using with IAM Roles

When running on EC2, ECS, or EKS, you can omit credentials to use IAM roles:

```bash
# No credentials needed - uses instance IAM role
oci-s3 push my-bucket/myapp:v1.0 --region us-east-1
```

## Troubleshooting

**Authentication errors**: Ensure AWS credentials are properly configured
```bash
aws configure
# or set environment variables
export AWS_ACCESS_KEY_ID=xxx
export AWS_SECRET_ACCESS_KEY=yyy
export AWS_REGION=us-east-1
```

**Docker daemon errors**: Ensure Docker is running
```bash
docker ps
```

## License

MIT License - see LICENSE file for details
