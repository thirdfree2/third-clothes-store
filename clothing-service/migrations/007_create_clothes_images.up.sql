CREATE TABLE clothes_images (
    id BIGSERIAL PRIMARY KEY,

    clothes_id BIGINT NOT NULL,
    bucket VARCHAR(100) NOT NULL,
    object_key TEXT NOT NULL,

    original_filename TEXT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_clothes_images_clothes
        FOREIGN KEY (clothes_id)
        REFERENCES clothes(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_clothes_images_clothes_id
    ON clothes_images(clothes_id);

CREATE UNIQUE INDEX idx_clothes_images_object_key
    ON clothes_images(object_key);