DROP INDEX IF EXISTS idx_clothes_color_id;

ALTER TABLE clothes
DROP CONSTRAINT IF EXISTS fk_clothes_color;

ALTER TABLE clothes
DROP COLUMN IF EXISTS color_id;