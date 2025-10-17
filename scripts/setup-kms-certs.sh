#!/bin/bash

# Setup script for Cosmian KMS certificates
# This script generates the necessary certificates and P12 file for KMS

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
COSMIAN_DIR="$PROJECT_ROOT/cosmian"

echo "=========================================="
echo "Cosmian KMS Certificate Setup"
echo "=========================================="

# Create cosmian directory if it doesn't exist
if [ ! -d "$COSMIAN_DIR" ]; then
    echo "Creating cosmian directory..."
    mkdir -p "$COSMIAN_DIR"
fi

cd "$COSMIAN_DIR"

# Check if certificates already exist
if [ -f "kms.server.p12" ] && [ -f "kms.crt" ] && [ -f "kms.key" ]; then
    echo ""
    echo "Certificates already exist in $COSMIAN_DIR"
    read -p "Do you want to regenerate them? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Using existing certificates."
        exit 0
    fi
    echo "Regenerating certificates..."
fi

echo ""
echo "Step 1: Creating private key and self-signed certificate..."
openssl req -newkey rsa:4096 -nodes -keyout kms.key -x509 -days 365 -out kms.crt -subj "/CN=kms.local"

if [ $? -eq 0 ]; then
    echo "✓ Private key and certificate created successfully"
else
    echo "✗ Failed to create private key and certificate"
    exit 1
fi

echo ""
echo "Step 2: Creating PKCS#12 (.p12) file..."
openssl pkcs12 -export -out kms.server.p12 -inkey kms.key -in kms.crt -password pass:password

if [ $? -eq 0 ]; then
    echo "✓ P12 file created successfully"
else
    echo "✗ Failed to create P12 file"
    exit 1
fi

echo ""
echo "=========================================="
echo "Certificate generation complete!"
echo "=========================================="
echo ""
echo "Generated files in $COSMIAN_DIR:"
echo "  - kms.key (Private key)"
echo "  - kms.crt (Certificate)"
echo "  - kms.server.p12 (PKCS#12 file)"
echo ""
echo "You can now run: docker compose up"
echo ""
