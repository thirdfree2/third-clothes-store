ALTER TABLE purchase_items
DROP CONSTRAINT IF EXISTS purchase_items_unit_price_check;

ALTER TABLE purchase_items
DROP COLUMN IF EXISTS unit_price;

ALTER TABLE purchases
DROP CONSTRAINT IF EXISTS purchases_total_amount_check;

ALTER TABLE purchases
DROP COLUMN IF EXISTS total_amount;
