CREATE TABLE IF NOT EXISTS clothes_categories (
    clothes_id bigint NOT NULL,
    category_id bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY (clothes_id, category_id),

    CONSTRAINT fk_clothes_categories_clothes
        FOREIGN KEY (clothes_id) REFERENCES clothes(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_clothes_categories_category
        FOREIGN KEY (category_id) REFERENCES categories(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_clothes_categories_clothes_id
ON clothes_categories(clothes_id);

CREATE INDEX IF NOT EXISTS idx_clothes_categories_category_id
ON clothes_categories(category_id);