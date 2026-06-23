-- 009_add_soft_delete_to_clothes_and_images.down.sql

DROP INDEX IF EXISTS idx_clothes_images_clothes_id_deleted_at;
DROP INDEX IF EXISTS idx_clothes_images_deleted_at;
DROP INDEX IF EXISTS idx_clothes_deleted_at;

ALTER TABLE clothes_images
DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE clothes
DROP COLUMN IF EXISTS deleted_at;