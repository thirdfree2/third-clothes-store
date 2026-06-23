ALTER TABLE clothes
ADD COLUMN IF NOT EXISTS color_id bigint;

ALTER TABLE clothes
ADD CONSTRAINT fk_clothes_color
FOREIGN KEY (color_id) REFERENCES colors(id);

CREATE INDEX IF NOT EXISTS idx_clothes_color_id
ON clothes(color_id);