CREATE TABLE IF NOT EXISTS admin_users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    role VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT admin_users_role_check
        CHECK (role IN ('ADMIN_VIEWER', 'ADMIN_EDITOR', 'ADMIN_OWNER')),

    CONSTRAINT admin_users_status_check
        CHECK (status IN ('ACTIVE', 'INACTIVE', 'BANNED'))
);

CREATE INDEX IF NOT EXISTS idx_admin_users_email
ON admin_users(email);

INSERT INTO admin_users (
    id,
    email,
    password_hash,
    role,
    status,
    created_at,
    updated_at
)
SELECT
    u.id,
    u.email,
    u.password_hash,
    ap.role,
    u.status,
    u.created_at,
    u.updated_at
FROM users u
JOIN admin_profiles ap ON ap.user_id = u.id
WHERE u.user_type = 'ADMIN'
ON CONFLICT (email) DO NOTHING;

SELECT SETVAL(
    PG_GET_SERIAL_SEQUENCE('admin_users', 'id'),
    COALESCE((SELECT MAX(id) FROM admin_users), 1),
    (SELECT COUNT(*) > 0 FROM admin_users)
);
