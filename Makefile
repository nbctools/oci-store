# Build configuration
OCI-S3_BINARY=oci-s3
GO_VERSION=1.21

.PHONY: build clean test install

# Build the binary
build:
	go build -o $(OCI-S3_BINARY) .

# Clean build artifacts
clean:
	rm -f $(OCI-S3_BINARY)

# Run tests
test:
	go test -v ./...

# Install binary to system
install: build
	sudo cp $(OCI-S3_BINARY) /usr/local/bin/

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(OCI-S3_BINARY)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(OCI-S3_BINARY)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o $(OCI-S3_BINARY)-windows-amd64.exe .