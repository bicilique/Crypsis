#!/usr/bin/env bash
set -euo pipefail

echo "🚀 Crypsis + Cosmian KMS Quick Start"
echo "====================================="
echo ""

# Helper: wait for a container command to succeed
wait_for() {
  local cmd="$1"
  local retries=${2:-30}
  local delay=${3:-2}
  local i=0
  until eval "$cmd" >/dev/null 2>&1; do
    i=$((i+1))
    if [ "$i" -ge "$retries" ]; then
      echo "Timeout waiting for command: $cmd"
      return 1
    fi
    sleep "$delay"
  done
  return 0
}

# 1) Ensure KMS certificates/keys exist (generate if missing)
if [ ! -f "cosmian/kms.server.p12" ] || [ ! -f "cosmian/kms.key" ] || [ ! -f "cosmian/kms.crt" ]; then
    echo "📝 Step 1: Generating KMS certificates and keys..."
    if [ -x "./scripts/setup-kms-certs.sh" ]; then
        ./scripts/setup-kms-certs.sh
    else
        echo "Error: ./scripts/setup-kms-certs.sh not found or not executable. Please create KMS certs under cosmian/ and rerun."
        exit 1
    fi
    echo ""
else
    echo "✓ KMS certificates and keys already exist"
    echo ""
fi

# 2) Ensure master key exists (create sample if missing)
if [ ! -f "resources/sample.key" ]; then
    echo "📝 Step 2: Generating sample master key at resources/sample.key..."
    mkdir -p resources
    # generate a 32-byte base64 key
    head -c 32 /dev/urandom | base64 > resources/sample.key
    echo "Generated resources/sample.key"
    echo ""
else
    echo "✓ Master key exists at resources/sample.key"
    echo ""
fi

# 3) Start services
echo "📝 Step 3: Starting all services (docker compose up -d)..."
docker compose up -d

echo ""

# 4) Wait for Postgres to be ready inside the 'db' service
echo "⏳ Waiting for Postgres to be ready..."
if ! wait_for "docker compose exec -T db pg_isready -U \"${POSTGRES_USER:-postgres}\"" 60 2; then
  echo "Postgres did not become ready in time. Check 'docker compose logs db'"
  docker compose logs --no-color db --tail=200
  exit 1
fi

echo "✓ Postgres is ready"

echo "\n🔁 Running one-time initialization tasks..."

# 5) Create MinIO buckets and users (run createbuckets service if defined)
if docker compose ps --services | grep -q "createbuckets"; then
  echo "🧰 Creating MinIO buckets and user (createbuckets)..."
  docker compose run --rm createbuckets || echo "Warning: createbuckets service failed — check logs"
  echo ""
fi

# 6) Initialize Hydra (run hydra-init if present)
if docker compose ps --services | grep -q "hydra-init"; then
  echo "🔐 Initializing Hydra (hydra-init)..."
  docker compose run --rm hydra-init || echo "Warning: hydra-init failed — check logs"
  echo ""
fi

# 7) Wait a bit and show service status
echo "⏳ Waiting for services to settle..."
sleep 5

echo "\n📊 Service Status:"
docker compose ps

echo ""
echo "✅ Setup complete!"
echo ""
echo "Services available at:"
echo "  - MinIO Console: http://localhost:9001"
echo "  - MinIO API: http://localhost:9000"
echo "  - PostgreSQL: postgresql://localhost:5432"
echo "  - Hydra Admin: http://localhost:4445"
echo "  - Hydra Public: http://localhost:4444"
echo "  - KMS (HTTPS): https://localhost:9998"
echo "  - Frontend: http://localhost:80"
echo ""
echo "To view logs:    docker compose logs -f"
echo "To stop:         docker compose down"
echo ""
