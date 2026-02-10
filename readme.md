# OCI-Store: Multi-Backend Docker Image Storage Manager

A CLI tool to store Docker/OCI images directly in cloud storage using Docker Distribution's production-grade storage drivers. Supports S3, Google Cloud Storage, and Azure Blob Storage.

## Features

- 🗄️ **Multiple Storage Backends**: S3, GCS, Azure with identical CLI experience
- 🔄 **Content-Addressable Storage**: Automatic deduplication and optimization
- 📦 **OCI Compliant**: Full compatibility with Docker and OCI image specifications
- ☁️ **Serverless**: No dedicated registry.


## Quick Start

### Install

```bash
go install github.com/nbctools/oci-store
```

### Build

```bash
go build -o oci-store .
```

## Usage

### S3 Storage

```bash
# Push to S3
oci-store s3 push --region us-east-1 my-bucket/myapp:v1.0

# Pull from S3
oci-store s3 pull --region us-east-1 my-bucket/myapp:v1.0

# With custom endpoint (S3-compatible storage)
oci-store s3 push --region us-east-1 --endpoint https://minio.local:9000 my-bucket/myapp:v1.0

# Push using an explicit image
oci-store s3 push --region us-east-1 my-bucket/myapp:latest --image mylocalapp:latest
```

For authentication and permissions see https://distribution.github.io/distribution/storage-drivers/s3/ 

### Google Cloud Storage

```bash
# Push to GCS
oci-store gcs push --project-id my-project my-bucket/myapp:v1.0

# Pull from GCS
oci-store gcs pull --project-id my-bucket/myapp:v1.0
```
For authentication and permissions see https://distribution.github.io/distribution/storage-drivers/gcs/


### Azure Blob Storage

```bash
# Push to Azure
oci-store azure push --account-name myaccount --account-key mykey my-container/myapp:v1.0

# Pull from Azure
oci-store azure pull --account-name myaccount --account-key mykey my-container/myapp:v1.0
```

For authentication and permission see https://distribution.github.io/distribution/storage-drivers/azure/

## Prerequisites

- Go 1.21 or later
- Docker installed and running
- Cloud Storage account with valid permissions see https://distribution.github.io/distribution/storage-drivers/


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


## Development

### Building

```bash
go build -o oci-store .
```

## License

This project uses the same license as Docker Distribution.
