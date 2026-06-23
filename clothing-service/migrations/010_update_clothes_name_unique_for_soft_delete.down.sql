DROP INDEX IF EXISTS idx_clothes_name_unique_active;

CREATE UNIQUE INDEX idx_clothes_name_unique
ON clothes (LOWER(name));