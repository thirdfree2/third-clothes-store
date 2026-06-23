ALTER TABLE purchases
ADD COLUMN IF NOT EXISTS total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0;

ALTER TABLE purchases
DROP CONSTRAINT IF EXISTS purchases_total_amount_check;

ALTER TABLE purchases
ADD CONSTRAINT purchases_total_amount_check
CHECK (total_amount >= 0);

ALTER TABLE purchase_items
ADD COLUMN IF NOT EXISTS unit_price NUMERIC(12, 2) NOT NULL DEFAULT 0;

ALTER TABLE purchase_items
DROP CONSTRAINT IF EXISTS purchase_items_unit_price_check;

ALTER TABLE purchase_items
ADD CONSTRAINT purchase_items_unit_price_check
CHECK (unit_price >= 0);
