CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,

    user_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_type_check
        CHECK (user_type IN ('ADMIN', 'CUSTOMER')),

    CONSTRAINT users_status_check
        CHECK (status IN ('ACTIVE', 'INACTIVE', 'BANNED'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
ON users (LOWER(email));

CREATE INDEX IF NOT EXISTS idx_users_user_type
ON users(user_type);

CREATE TABLE IF NOT EXISTS admin_profiles (
    user_id BIGINT PRIMARY KEY,

    role VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_admin_profiles_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT admin_profiles_role_check
        CHECK (role IN ('ADMIN_VIEWER', 'ADMIN_EDITOR', 'ADMIN_OWNER'))
);

INSERT INTO users (
    id,
    email,
    password_hash,
    user_type,
    status,
    created_at,
    updated_at
)
SELECT
    id,
    LOWER(email),
    password_hash,
    'ADMIN',
    status,
    created_at,
    updated_at
FROM admin_users
ON CONFLICT (LOWER(email)) DO NOTHING;

INSERT INTO admin_profiles (
    user_id,
    role,
    created_at,
    updated_at
)
SELECT
    au.id,
    au.role,
    au.created_at,
    au.updated_at
FROM admin_users au
JOIN users u ON u.id = au.id
WHERE u.user_type = 'ADMIN'
ON CONFLICT (user_id) DO NOTHING;

SELECT SETVAL(
    PG_GET_SERIAL_SEQUENCE('users', 'id'),
    COALESCE((SELECT MAX(id) FROM users), 1),
    (SELECT COUNT(*) > 0 FROM users)
);
