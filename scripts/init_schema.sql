-- 1. Admins table (independent)
CREATE TABLE admins (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    secret VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    salt VARCHAR(255) NOT NULL
);
CREATE INDEX idx_admins_deleted_at ON admins (deleted_at);

-- 2. Apps table (independent)
CREATE TABLE apps (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL,
    uri TEXT,
    redirect_uri TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_apps_client_id ON apps (client_id);
CREATE INDEX idx_apps_is_active ON apps (is_active);
CREATE INDEX idx_apps_deleted_at ON apps (deleted_at);

-- 3. Files table (independent, referenced by metadata)
CREATE TABLE files (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    app_id VARCHAR(36),
    user_id VARCHAR(36),
    mime_type VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    bucket_name VARCHAR(255),
    location TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_files_name ON files (name);
CREATE INDEX idx_files_app_id ON files (app_id);
CREATE INDEX idx_files_user_id ON files (user_id);
CREATE INDEX idx_files_deleted_at ON files (deleted_at);

-- 4. FileLogs table (can reference files via file_id, if needed)
CREATE TABLE file_logs (
    id SERIAL PRIMARY KEY,
    actor_id TEXT NOT NULL,
    actor_type TEXT NOT NULL CHECK (actor_type IN ('user', 'client', 'system', 'admin')),
    file_id UUID NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('upload', 'download', 'update', 'delete', 'recover', 'encrypt', 'decrypt', 're-key')),
    timestamp TIMESTAMPTZ DEFAULT now(),
    ip INET,
    user_agent TEXT,
    metadata JSONB
);
CREATE INDEX idx_file_logs_file_id ON file_logs (file_id);

-- 5. Metadata table (depends on files)
CREATE TABLE metadata (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    file_id VARCHAR(36) NOT NULL,
    hash VARCHAR(256) NOT NULL,
    enc_hash VARCHAR(256),
    key_uid VARCHAR(256),
    enc_key TEXT NOT NULL,
    key_algo VARCHAR(64) NOT NULL,
    version_id VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_metadata_file FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);
CREATE INDEX idx_metadata_file_id ON metadata (file_id);
CREATE INDEX idx_metadata_enc_hash ON metadata (enc_hash);
CREATE INDEX idx_metadata_key_uid ON metadata (key_uid);
CREATE INDEX idx_metadata_deleted_at ON metadata (deleted_at);
