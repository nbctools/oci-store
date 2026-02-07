#!/bin/bash

# Quick Multi-Backend Test Script
# Quickly validates all three storage backends

set -e

echo "🧪 Quick Multi-Backend Validation"

# Check if OCI-Store binary exists
if [ ! -f "oci-store" ]; then
    echo "❌ Building oci-store..."
    go build -o oci-store .
fi

# Test help commands
echo ""
echo "📋 Testing CLI structure..."
echo "=== Main Help ==="
./oci-store --help | head -10
echo ""
echo "=== S3 Commands ==="
./oci-store s3 --help | head -10
echo ""
echo "=== GCS Commands ==="
./oci-store gcs --help | head -10
echo ""
echo "=== Azure Commands ==="
./oci-store azure --help | head -10

echo ""
echo "🎯 Testing validation..."

# Test S3 validation
echo "=== S3 Validation (should fail) ==="
./oci-store s3 push test-bucket/app:v1.0 2>&1 | grep -q "requires region" && echo "✅ S3 validation working" || echo "❌ S3 validation failed"

# Test Azure validation
echo "=== Azure Validation (should fail) ==="
./oci-store azure push test-container/app:v1.0 2>&1 | grep -q "requires account name" && echo "✅ Azure validation working" || echo "❌ Azure validation failed"

echo ""
echo "🎉 CLI validation completed!"
echo ""
echo "To run full demo: ./demo/multi-backend-demo.sh"
