CREATE TABLE IF NOT EXISTS clothes (
    id bigserial PRIMARY KEY,
    name varchar(180) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_clothes_name ON clothes (name);
