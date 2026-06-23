CREATE TABLE IF NOT EXISTS categories (
    id bigserial PRIMARY KEY,
    name varchar(120) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name
ON categories (name);