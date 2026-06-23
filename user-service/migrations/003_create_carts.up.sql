CREATE TABLE IF NOT EXISTS carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT carts_status_check
        CHECK (status IN ('ACTIVE', 'CHECKED_OUT'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_carts_active_user
ON carts(user_id)
WHERE status = 'ACTIVE';

CREATE INDEX IF NOT EXISTS idx_carts_user_id
ON carts(user_id);

CREATE TABLE IF NOT EXISTS cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL,
    clothes_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    size VARCHAR(20) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_cart_items_cart
        FOREIGN KEY (cart_id)
        REFERENCES carts(id)
        ON DELETE CASCADE,

    CONSTRAINT cart_items_quantity_check
        CHECK (quantity > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cart_items_unique_item
ON cart_items(cart_id, clothes_id, size);

CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id
ON cart_items(cart_id);
