# oci-store

Store and retrieve OCI artifacts directly from object storage (S3, GCS, Azure)
— without running a container registry.

Perfect for CI/CD pipelines, air-gapped environments, backups, and low-ops setups.

## Why oci-store?

Running a full OCI registry (Harbor, Docker Registry, etc.) can be:
- Operationally heavy
- Overkill for CI artifacts or internal use
- Hard to justify for ephemeral or air-gapped environments

`oci-store` lets you:
- Store OCI artifacts directly in object storage
- Avoid running and maintaining a registry service
- Reuse cheap, durable cloud storage
- Integrate cleanly with existing OCI tooling


## Quick start

Install 

```bash
go install github.com/nbctools/oci-store
```

Push an OCI artifact to S3:

```bash
oci-store s3 push \
    --region us-east-1 \
    --image myimage:latest \
    my-bucket/myapp:v1.0
```

Pull it back:

``` bash
oci-store s3 pull \
    --region us-east-1 \
    my-bucket/myapp:v1.0
```


## How It Works

`oci-store` acts as a bridge between local OCI images and object storage:

```
┌─────────────────┐
│  Local Docker   │
│     Daemon      │
└────────┬────────┘
         │
         │ oci-store reads/writes
         │ via Docker API
         ▼
┌─────────────────┐      ┌──────────────────┐
│   oci-store     │◄────►│ Object Storage   │
│   (This CLI)    │      │ (S3/GCS/Azure)   │
└─────────────────┘      └──────────────────┘
                               │
                               ▼
                         Registry v2 Layout
                         (blobs + metadata)
```

**Push workflow:**
1. Reads OCI image from local Docker daemon
2. Extracts layers, manifests, and configs
3. Uploads them to object storage using registry v2 layout
4. Content-addressable storage ensures deduplication

**Pull workflow:**
1. Downloads manifests and layers from object storage
2. Reconstructs OCI image format
3. Loads image into local Docker daemon

No persistent registry service — `oci-store` uses an ephemeral registry process only during push/pull operations, then tears it down.

## Usage

### AWS S3 Storage

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

- Docker daemon installed and running
- Cloud Storage account with valid permissions see https://distribution.github.io/distribution/storage-drivers/

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


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
