CREATE TABLE admin_users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    role VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT admin_users_role_check
        CHECK (role IN ('ADMIN_VIEWER', 'ADMIN_EDITOR')),

    CONSTRAINT admin_users_status_check
        CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE INDEX idx_admin_users_email
ON admin_users(email);