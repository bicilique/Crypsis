#!/bin/bash

echo "=== Crypsis KEK Diagnostic Tool ==="
echo ""

# Check environment variables
echo "Checking environment variables..."
echo "KMS_ENABLE: ${KMS_ENABLE:-not set}"
echo "KMS_KEY_UID: ${KMS_KEY_UID:-not set}"
echo "MKEY_PATH: ${MKEY_PATH:-not set}"
echo ""

# Check if KMS is enabled
if [ "${KMS_ENABLE}" = "true" ]; then
    echo "✓ KMS is ENABLED"
    echo "  KEK will be loaded from KMS with UID: $KMS_KEY_UID"
    echo ""
    echo "IMPORTANT: Make sure the KMS_KEY_UID points to a valid AES-256 key (32 bytes)"
    echo ""
else
    echo "✗ KMS is DISABLED"
    echo "  KEK will be loaded from file: ${MKEY_PATH:-resources/sample.key}"
    echo ""
    
    # Check if the file exists
    KEK_FILE="${MKEY_PATH:-resources/sample.key}"
    if [ -f "$KEK_FILE" ]; then
        echo "✓ KEK file exists: $KEK_FILE"
        FILE_SIZE=$(wc -c < "$KEK_FILE" | tr -d ' ')
        echo "  File size: $FILE_SIZE bytes"
        
        # Check if it's a valid Tink keyset (should be around 100+ bytes)
        if [ "$FILE_SIZE" -lt 50 ]; then
            echo "  ⚠️  WARNING: File is too small to be a valid Tink keyset!"
            echo "  Expected size: ~100-150 bytes for a binary Tink keyset"
        elif [ "$FILE_SIZE" -eq 32 ]; then
            echo "  ⚠️  WARNING: File is exactly 32 bytes - this looks like a raw AES key, not a Tink keyset!"
            echo "  You need to convert it to a Tink keyset format"
        elif [ "$FILE_SIZE" -eq 44 ] || [ "$FILE_SIZE" -eq 64 ]; then
            echo "  ⚠️  WARNING: File size suggests it's a base64-encoded raw key, not a Tink keyset!"
            echo "  You need to convert it to a Tink keyset format"
        else
            echo "  ✓ File size looks reasonable for a Tink keyset"
        fi
        
        # Try to read first few bytes to check if it's binary or text
        FIRST_BYTES=$(head -c 10 "$KEK_FILE" | xxd -p | tr -d '\n')
        echo "  First bytes (hex): $FIRST_BYTES"
        
        # Check if it starts with Tink keyset magic bytes
        # Tink keysets typically start with 08 or 0a (protobuf wire types)
        if [[ "$FIRST_BYTES" =~ ^(08|0a).* ]]; then
            echo "  ✓ File starts with protobuf-like bytes (good sign for Tink keyset)"
        else
            echo "  ⚠️  File doesn't start with expected protobuf bytes"
        fi
    else
        echo "✗ KEK file NOT FOUND: $KEK_FILE"
        echo "  Please create the file or update MKEY_PATH in your .env"
    fi
fi

echo ""
echo "=== Recommendations ==="
echo ""
echo "1. If using KMS (KMS_ENABLE=true):"
echo "   - Make sure KMS_KEY_UID points to a master key (KEK), not a data key"
echo "   - The KEK from KMS will be automatically converted to Tink format"
echo ""
echo "2. If NOT using KMS (KMS_ENABLE=false):"
echo "   - The KEK file must be a binary Tink keyset (not a raw hex or base64 key)"
echo "   - Use the Go test 'TestHexKeyFromKMS' to verify key conversion"
echo "   - To regenerate a valid KEK, run: go run scripts/generate_kek.go"
echo ""
echo "3. Common issues:"
echo "   - KEK is a raw hex string (not converted to Tink keyset)"
echo "   - KEK is base64-encoded raw bytes (not a Tink keyset)"
echo "   - Wrong KMS_KEY_UID (pointing to DEK instead of KEK)"
echo ""
