CREATE TABLE IF NOT EXISTS customer_profiles (
    user_id BIGINT PRIMARY KEY,

    first_name VARCHAR(100) NULL,
    last_name VARCHAR(100) NULL,
    phone VARCHAR(30) NULL,
    date_of_birth DATE NULL,
    marketing_opt_in BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_customer_profiles_phone
ON customer_profiles(phone)
WHERE phone IS NOT NULL;
