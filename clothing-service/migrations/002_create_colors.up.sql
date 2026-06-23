CREATE TABLE IF NOT EXISTS colors (
    id bigserial PRIMARY KEY,
    name varchar(120) NOT NULL,
    hex_code varchar(20),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_colors_name
ON colors (name);