#!/bin/bash
set -euo pipefail

CERT_DIR="cosmian"
mkdir -p "$CERT_DIR"

# Generate private key and certificate
openssl req -x509 \
    -newkey rsa:4096 \
    -keyout "$CERT_DIR/kms.key" \
    -out "$CERT_DIR/kms.crt" \
    -days 365 \
    -nodes \
    -subj "/CN=crypsis-kms"

# Create PKCS12 file
openssl pkcs12 -export \
    -in "$CERT_DIR/kms.crt" \
    -inkey "$CERT_DIR/kms.key" \
    -out "$CERT_DIR/kms.server.p12" \
    -name "kms-cert" \
    -password pass:password

echo "✅ Generated KMS certificates in $CERT_DIR:"
ls -l "$CERT_DIR"
