#!/bin/bash
set -euo pipefail

echo "Running database initialization script..."

# Export PostgreSQL password for non-interactive authentication
export PGPASSWORD="${POSTGRES_PASSWORD:-}"

# Verify required PostgreSQL environment variables
if [[ -z "${POSTGRES_USER:-}" || -z "${POSTGRES_PASSWORD:-}" ]]; then
  echo "Error: POSTGRES_USER and POSTGRES_PASSWORD must be set in the environment."
  exit 1
fi

# Function to check if a database exists
database_exists() {
  [[ "$(psql -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_database WHERE datname = '$1'")" == "1" ]]
}

# Function to check if a user exists
user_exists() {
  [[ "$(psql -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_roles WHERE rolname = '$1'")" == "1" ]]
}

# Function to create database and assign user with permissions
create_db_and_user() {
  local db_name="$1"
  local db_user="$2"
  local db_password="$3"

  if ! database_exists "$db_name"; then
    echo "Creating database: $db_name"
    psql -U "$POSTGRES_USER" -c "CREATE DATABASE \"$db_name\";"
  else
    echo "Database $db_name already exists."
  fi

  if ! user_exists "$db_user"; then
    echo "Creating user: $db_user"
    psql -U "$POSTGRES_USER" -c "CREATE USER \"$db_user\" WITH ENCRYPTED PASSWORD '$db_password';"
  else
    echo "User $db_user already exists."
  fi

  echo "Granting privileges to user $db_user on database $db_name"
  psql -U "$POSTGRES_USER" -c "GRANT ALL PRIVILEGES ON DATABASE \"$db_name\" TO \"$db_user\";"
  psql -U "$POSTGRES_USER" -d "$db_name" -c "GRANT USAGE, CREATE ON SCHEMA public TO \"$db_user\";"
  psql -U "$POSTGRES_USER" -d "$db_name" -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO \"$db_user\";"
  psql -U "$POSTGRES_USER" -d "$db_name" -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO \"$db_user\";"
  psql -U "$POSTGRES_USER" -d "$db_name" -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO \"$db_user\";"
}

# Create and configure databases/users
create_db_and_user "$TEST_DEV_DB" "$TEST_DEV_USER" "$TEST_DEV_PASSWORD"
create_db_and_user "$HYDRA_DB" "$HYDRA_DB_USER" "$HYDRA_DB_PASSWORD"

# Create KMS database
echo "Creating KMS database..."
create_db_and_user "kms_db" "$POSTGRES_USER" "$POSTGRES_PASSWORD"

# Load schema file
SCHEMA_PATH="$(dirname "$0")/init_schema.sql"

if [[ ! -f "$SCHEMA_PATH" ]]; then
  echo "Error: Schema file not found at $SCHEMA_PATH"
  exit 1
fi

echo "Applying schema to $TEST_DEV_DB..."
psql -U "$POSTGRES_USER" -d "$TEST_DEV_DB" -f "$SCHEMA_PATH"
echo "Schema applied successfully!"

# Grant dev_user access and change ownership (fix for automigrate)
echo "Granting access and changing ownership to $TEST_DEV_USER..."
psql -U "$POSTGRES_USER" -d "$TEST_DEV_DB" <<EOF
-- Grant access on all existing objects
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "$TEST_DEV_USER";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "$TEST_DEV_USER";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO "$TEST_DEV_USER";

-- Change ownership of all tables
DO \$\$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT tablename FROM pg_tables WHERE schemaname = 'public' LOOP
        EXECUTE 'ALTER TABLE public.' || quote_ident(r.tablename) || ' OWNER TO "$TEST_DEV_USER";';
    END LOOP;

    FOR r IN SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public' LOOP
        EXECUTE 'ALTER SEQUENCE public.' || quote_ident(r.sequence_name) || ' OWNER TO "$TEST_DEV_USER";';
    END LOOP;

    FOR r IN SELECT routine_name FROM information_schema.routines WHERE routine_schema = 'public' LOOP
        BEGIN
          EXECUTE 'ALTER FUNCTION public.' || quote_ident(r.routine_name) || '() OWNER TO "$TEST_DEV_USER";';
        EXCEPTION WHEN OTHERS THEN
          -- Skip functions with unsupported signatures
          CONTINUE;
        END;
    END LOOP;
END
\$\$;
EOF

# Insert default admin user
echo "Creating default admin user..."

REQUIRED_VARS=(DEFAULT_ADMIN_USERNAME DEFAULT_ADMIN_CLIENT_ID DEFAULT_ADMIN_SECRET DEFAULT_ADMIN_SALT)
for var in "${REQUIRED_VARS[@]}"; do
  if [[ -z "${!var:-}" ]]; then
    echo "Error: Environment variable $var must be set for admin creation."
    exit 1
  fi
done

if psql -U "$POSTGRES_USER" -d "$TEST_DEV_DB" -tAc "SELECT to_regclass('public.admins')" | grep -q 'admins'; then
  echo "Ensuring pgcrypto extension is available..."
  psql -U "$POSTGRES_USER" -d "$TEST_DEV_DB" -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"

  echo "Inserting default admin into 'admins' table (if not exists)..."
  psql -U "$POSTGRES_USER" -d "$TEST_DEV_DB" <<EOF
INSERT INTO admins (id, username, client_id, secret, created_at, updated_at, salt)
VALUES (
  gen_random_uuid(),
  '${DEFAULT_ADMIN_USERNAME}',
  '${DEFAULT_ADMIN_CLIENT_ID}',
  '${DEFAULT_ADMIN_SECRET}',
  now(),
  now(),
  '${DEFAULT_ADMIN_SALT}'
)
ON CONFLICT (username) DO NOTHING;
EOF

  echo "Default admin user created (if not already present)."
else
  echo "Warning: Table 'admins' does not exist in database '$TEST_DEV_DB'. Skipping admin insert."
fi

echo "Database initialization complete."
