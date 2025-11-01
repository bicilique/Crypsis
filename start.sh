#!/bin/bash
# Quick start script for Crypsis project with KMS

echo "🚀 Crypsis + Cosmian KMS Quick Start"
echo "====================================="
echo ""

# Check if certificates exist
if [ ! -f "cosmian/kms.server.p12" ]; then
    echo "📝 Step 1: Setting up KMS certificates..."
    ./scripts/setup-kms-certs.sh
    echo ""
else
    echo "✓ KMS certificates already exist"
    echo ""
fi

echo "📝 Step 2: Starting all services..."
docker compose up -d

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 5

echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "✅ Setup complete!"
echo ""
echo "Services available at:"
echo "  - MinIO:        http://localhost:9000"
echo "  - MinIO Console: http://localhost:9001"
echo "  - PostgreSQL:    postgresql://localhost:5432"
echo "  - Hydra Admin:   http://localhost:4445"
echo "  - Hydra Public:  http://localhost:4444"
echo "  - KMS:          https://localhost:9998"
echo ""
echo "To view logs:    docker compose logs -f"
echo "To stop:         docker compose down"
echo ""
