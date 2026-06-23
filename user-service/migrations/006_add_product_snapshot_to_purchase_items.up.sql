ALTER TABLE purchase_items
ADD COLUMN IF NOT EXISTS product_name VARCHAR(180) NOT NULL DEFAULT '';

ALTER TABLE purchase_items
ADD COLUMN IF NOT EXISTS product_image_url TEXT;
