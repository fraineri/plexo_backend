CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    status VARCHAR(10) NOT NULL CHECK (status IN ('ALLOW', 'DENY')),
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    PRIMARY KEY (role_id, permission_id)
);
