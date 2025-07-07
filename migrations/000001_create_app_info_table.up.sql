CREATE TABLE IF NOT EXISTS app_info (
    name VARCHAR(255) PRIMARY KEY,
    version VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    deleted_at BIGINT
);

INSERT INTO app_info (name, version, status) VALUES ('plexo', '0.1.0', 'ENABLED');
