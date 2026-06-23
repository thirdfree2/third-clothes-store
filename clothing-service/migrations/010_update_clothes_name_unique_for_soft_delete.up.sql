ALTER TABLE clothes
DROP CONSTRAINT IF EXISTS clothes_name_key;

ALTER TABLE clothes
DROP CONSTRAINT IF EXISTS idx_clothes_name;

DROP INDEX IF EXISTS idx_clothes_name;
DROP INDEX IF EXISTS idx_clothes_name_unique;
DROP INDEX IF EXISTS idx_clothes_name_unique_active;

CREATE UNIQUE INDEX idx_clothes_name_unique_active
ON clothes (LOWER(name))
WHERE deleted_at IS NULL;