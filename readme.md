# OCI-Store: Multi-Backend Docker Image Storage

A CLI tool to store Docker/OCI images directly in cloud storage using Docker Distribution's production-grade storage drivers. Supports S3, Google Cloud Storage, and Azure Blob Storage.

## Features

- 🗄️ **Multiple Storage Backends**: S3, GCS, Azure with identical CLI experience
- 🚀 **Local Development**: Emulate all backends locally with LocalStack and Azurite
- 🔄 **Content-Addressable Storage**: Automatic deduplication and optimization
- 📦 **OCI Compliant**: Full compatibility with Docker and OCI image specifications
- 🛡️ **Production Ready**: Uses battle-tested Docker Distribution storage drivers

## Quick Start

### Build

```bash
go build -o oci-store .
```

### Quick Test

```bash
# Validate CLI structure and validation
./demo/quick-test.sh
```

### Multi-Backend Demo

```bash
# Start local emulators and test all backends
./demo/multi-backend-demo.sh
```

## Usage

### S3 Storage

```bash
# Push to S3
./oci-store s3 push --region us-east-1 my-bucket/myapp:v1.0

# Pull from S3
./oci-store s3 pull --region us-east-1 my-bucket/myapp:v1.0

# With custom endpoint (S3-compatible storage)
./oci-store s3 push --region us-east-1 --endpoint https://minio.local:9000 my-bucket/myapp:v1.0
```

**Environment Variables:**
- `AWS_REGION`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### Google Cloud Storage

```bash
# Push to GCS
./oci-store gcs push --project-id my-project my-bucket/myapp:v1.0

# Pull from GCS
./oci-store gcs pull --project-id my-bucket/myapp:v1.0
```

**Environment Variables:**
- `GOOGLE_CLOUD_PROJECT`

### Azure Blob Storage

```bash
# Push to Azure
./oci-store azure push --account-name myaccount --account-key mykey my-container/myapp:v1.0

# Pull from Azure
./oci-store azure pull --account-name myaccount --account-key mykey my-container/myapp:v1.0
```

**Environment Variables:**
- `AZURE_STORAGE_ACCOUNT`
- `AZURE_STORAGE_KEY`

Or build from source:

```bash
git clone https://github.com/nbctools/oci-store
cd oci-store
go build -o oci-store
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
oci-store push my-bucket/myapp:v1.0

# Push with explicit local image name
oci-store push my-bucket/myapp:v1.0 --image localhost/myapp:latest

# Push to a specific region
oci-store push my-bucket/myapp:v1.0 --region us-west-2

# Push to S3-compatible storage (e.g., MinIO)
oci-store push my-bucket/myapp:v1.0 --endpoint http://localhost:9000

# Use a root directory in the bucket
oci-store push my-bucket/myapp:v1.0 --root-dir /registry
```

### Pull an image from S3

```bash
# Pull an image from S3 and load into Docker
oci-store pull my-bucket/myapp:v1.0
```

## Local Development

### Using Emulators

The project includes Docker Compose configurations for local development:

```bash
cd demo
docker-compose up -d
```

This starts:
- **LocalStack** (port 4566): Emulates S3, GCS, and Azure
- **Azurite** (port 10000): Dedicated Azure Blob Storage emulator

### Local S3 Testing

```bash
# Create bucket
aws --endpoint-url http://localhost:4566 s3 mb s3://test-bucket

# Push image
AWS_REGION=us-east-1 ./oci-store s3 push --endpoint http://localhost:4566 test-bucket/myapp:v1.0
```

### Local GCS Testing

```bash
# Create bucket (via LocalStack)
aws --endpoint-url http://localhost:4566 s3api create-bucket --bucket test-gcs-bucket

# Push image
GOOGLE_CLOUD_PROJECT=demo ./oci-store gcs push --endpoint http://localhost:4566 test-gcs-bucket/myapp:v1.0
```

### Local Azure Testing

```bash
# Create container (via Azurite)
curl -X PUT http://localhost:10000/devstoreaccount1/test-container?restype=container \
  -H "x-ms-version:2019-12-12" \
  -H "x-ms-date:$(date -u '+%a, %d %b %Y %H:%M:%S GMT')" \
  -H "Authorization:SharedKey devstoreaccount1:Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==" \
  -H "Content-Length:0"

# Push image
AZURE_STORAGE_ACCOUNT=devstoreaccount1 AZURE_STORAGE_KEY=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw== \
./oci-store azure push --account-name devstoreaccount1 --account-key Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw== test-container/myapp:v1.0
```

## CLI Reference

```
oci-store [command] [subcommand] [flags]

Commands:
  azure       Azure Blob Storage operations
  gcs         Google Cloud Storage operations
  s3          S3 storage operations

S3 Flags:
  --region            AWS region
  --endpoint          S3-compatible endpoint (optional)
  --access-key        AWS access key
  --secret-key        AWS secret key
  --root-dir          Root directory in bucket (optional)

GCS Flags:
  --project-id        GCS project ID
  --root-dir          Root directory in bucket (optional)

Azure Flags:
  --account-name      Storage account name
  --account-key       Storage account key
  --root-dir          Root directory in container (optional)

Global Flags:
  --verbose           Verbose output
```

## Storage Layout

All backends use the same Docker Distribution registry layout:

```
<storage>://<bucket>/<registry-v2-layout>
├── docker/
│   └── registry/
│       └── v2/
│           ├── blobs/        # Layer data (deduplicated)
│           └── repositories/ # Image metadata
```

This layout provides:
- **Content-addressable storage**: Blobs stored by digest for deduplication
- **Repository isolation**: Each image has its own namespace
- **Atomic operations**: Tags are symbolic links for atomic updates
- **Garbage collection support**: Unused blobs can be identified and removed

## Architecture

The tool uses a modular storage backend architecture:

```
oci-store
├── storage.go        # Storage interface and common logic
├── backends.go       # S3, GCS, Azure backend implementations
├── s3_commands.go   # S3-specific CLI commands
├── gcs_commands.go   # GCS-specific CLI commands
├── azure_commands.go # Azure-specific CLI commands
└── registry.go       # Docker Distribution registry management
```

Each storage backend implements the `StorageBackend` interface:
- `ParseRef()` - Parse storage-specific references
- `GetStorageConfig()` - Generate Docker Distribution storage configuration
- `ValidateConfig()` - Validate storage-specific configuration

## How It Works

This tool leverages the storage architecture from `distribution/distribution` (the official Docker Registry implementation):

1. **Storage Driver**: Uses production-grade storage drivers that power Docker Registry
2. **Registry Storage**: Implements the full registry storage API for proper blob and manifest management
3. **Content Addressable**: All blobs are stored by their SHA256 digest
4. **Optimizations**:
   - Multipart uploads for large layers (configurable chunk size)
   - Parallel multipart copy operations
   - Blob deduplication across all images
   - Transfer acceleration support (S3)
   - CDN integration capability

**Push Process**:
1. Retrieves image from local Docker daemon
2. Uses registry storage API to upload blobs (with automatic deduplication)
3. Uploads manifest and creates tag reference
4. Storage driver handles multipart uploads for large blobs

**Pull Process**:
1. Resolves tag to manifest digest via registry storage
2. Downloads manifest and parses layer references
3. Downloads all blobs using storage driver (with verification)
4. Reconstructs image and loads into Docker daemon

## Development

### Building

```bash
go build -o oci-store .
```

### Testing

```bash
# Quick validation
./demo/quick-test.sh

# Full multi-backend demo
./demo/multi-backend-demo.sh
```

### Dependencies

The project uses:
- `github.com/distribution/distribution/v3` - Core registry implementation
- `github.com/google/go-containerregistry` - OCI image handling
- `github.com/spf13/cobra` - CLI framework

## Permissions

### S3

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

### GCS

Your Google Cloud service account needs:
- `storage.buckets.get`
- `storage.objects.get`
- `storage.objects.create`
- `storage.objects.list`

### Azure

Your Azure storage account needs:
- Read and write access to blob containers
- List and create container permissions

## Benefits

- **Battle-tested**: Same storage backend used by Docker Registry in production
- **Deduplication**: Identical layers shared across images are stored only once
- **Multi-Cloud**: Consistent interface across all major cloud providers
- **Integrity**: Content hash ensures data integrity at every layer
- **Atomic Updates**: Tag updates are atomic operations
- **Efficient Transfer**: Only missing layers need to be uploaded/downloaded
- **Storage-Agnostic**: Works with S3-compatible storage, GCS, and Azure

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with `./demo/quick-test.sh`
5. Submit a pull request

## License

This project uses the same license as Docker Distribution.
