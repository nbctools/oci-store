# OCI-S3
A simple CLI tool written in Go that upload and downloads docker images from a S3 bucket. 
Supports the following features:
- Serverless: No dependency on a long running docker registry
- Optimised storage: Layers are deduplicated

It uses https://github.com/distribution/distribution/blob/main/registry/registry.go under the hood.

## Installation

### From Source
```bash
git clone <repository-url>
cd oci-s3
make install
```

### Build from Source
```bash
git clone <repository-url>
cd oci-s3
make build
```

## Usage

Push an image:

``` shell
oci-s3 push <s3-bucket>/<image>:<tag>
```

Pull an image:

``` shell
oci-s3 pull <s3-bucket>/<image>:<tag>
```

## Environment Variables

The tool requires AWS credentials to be available via environment variables:

- `AWS_ACCESS_KEY_ID`: Your AWS access key
- `AWS_SECRET_ACCESS_KEY`: Your AWS secret key
- `AWS_REGION`: AWS region (defaults to us-east-1)

## Examples

```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-west-2

# Push an image to S3 bucket "my-registry"
oci-s3 push my-registry/myapp:v1.0

# Pull an image from S3 bucket "my-registry"
oci-s3 pull my-registry/myapp:v1.0
```

## Development

### Building
```bash
make build
```

### Testing
```bash
make test
```

### Cross-compilation
```bash
make build-all
```

## Architecture

The tool follows OCI (Open Container Initiative) standards for image storage and uses S3 as the backend storage layer. Images are stored with proper layering and deduplication:

- `/<image>/<tag>/manifest.json` - OCI manifest file
- `/<image>/<tag>/layers/` - Individual layer files
- `/<image>/<tag>/config.json` - Image configuration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request
