CREATE TABLE IF NOT EXISTS app_info (
    name VARCHAR(255) PRIMARY KEY,
    version VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO app_info (name, version, status) VALUES ('plexo', '0.1.0', 'ENABLED');
