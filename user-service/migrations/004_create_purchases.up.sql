CREATE TABLE IF NOT EXISTS purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT purchases_total_amount_check
        CHECK (total_amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_purchases_user_id
ON purchases(user_id);

CREATE TABLE IF NOT EXISTS purchase_items (
    id BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    clothes_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    unit_price NUMERIC(12, 2) NOT NULL DEFAULT 0,
    size VARCHAR(20) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_purchase_items_purchase
        FOREIGN KEY (purchase_id)
        REFERENCES purchases(id)
        ON DELETE CASCADE,

    CONSTRAINT purchase_items_quantity_check
        CHECK (quantity > 0),

    CONSTRAINT purchase_items_unit_price_check
        CHECK (unit_price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_purchase_items_purchase_id
ON purchase_items(purchase_id);
