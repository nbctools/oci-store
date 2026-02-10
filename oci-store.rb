# Homebrew Formula for OCI-Store
# This file should be submitted to homebrew-core or used in a custom tap

class OciStore < Formula
  desc "Multi-backend Docker image storage CLI tool"
  homepage "https://github.com/nbctools/oci-store"
  url "https://github.com/nbctools/oci-store/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "sha256_placeholder"  # This will be updated by the release workflow
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    # Test basic CLI functionality
    assert_match "Multi-Backend Docker Image Storage", shell_output("#{bin}/oci-store --help")
    
    # Test subcommands exist
    assert_match "S3 storage operations", shell_output("#{bin}/oci-store s3 --help")
    assert_match "Google Cloud Storage operations", shell_output("#{bin}/oci-store gcs --help")
    assert_match "Azure Blob Storage operations", shell_output("#{bin}/oci-store azure --help")
    
    # Test validation works
    assert_match "requires region", shell_output("#{bin}/oci-store s3 push test-bucket/app:v1.0 2>&1", 1)
    assert_match "requires project ID", shell_output("#{bin}/oci-store gcs push test-bucket/app:v1.0 2>&1", 1)
    assert_match "requires account name", shell_output("#{bin}/oci-store azure push test-container/app:v1.0 2>&1", 1)
  end
end