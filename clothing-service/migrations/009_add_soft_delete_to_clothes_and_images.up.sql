-- 009_add_soft_delete_to_clothes_and_images.up.sql

ALTER TABLE clothes
ADD COLUMN deleted_at TIMESTAMPTZ NULL;

ALTER TABLE clothes_images
ADD COLUMN deleted_at TIMESTAMPTZ NULL;

CREATE INDEX idx_clothes_deleted_at
    ON clothes(deleted_at);

CREATE INDEX idx_clothes_images_deleted_at
    ON clothes_images(deleted_at);

CREATE INDEX idx_clothes_images_clothes_id_deleted_at
    ON clothes_images(clothes_id, deleted_at);