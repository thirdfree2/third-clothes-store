INSERT INTO colors (name, hex_code)
VALUES
    ('Black', '#000000'),
    ('White', '#FFFFFF'),
    ('Red', '#FF0000'),
    ('Blue', '#0000FF'),
    ('Green', '#00AA00')
ON CONFLICT (name) DO UPDATE
SET
    hex_code = EXCLUDED.hex_code,
    updated_at = now();

INSERT INTO categories (name)
VALUES
    ('Tops'),
    ('Pants'),
    ('Outerwear'),
    ('Shoes'),
    ('Accessories'),
    ('Casual'),
    ('Sport')
ON CONFLICT (name) DO NOTHING;

INSERT INTO clothes (name, price, color_id)
VALUES
    (
        'Basic T-Shirt',
        390.00,
        (SELECT id FROM colors WHERE name = 'Black')
    ),
    (
        'Blue Jeans',
        1290.00,
        (SELECT id FROM colors WHERE name = 'Blue')
    ),
    (
        'Red Hoodie',
        990.00,
        (SELECT id FROM colors WHERE name = 'Red')
    ),
    (
        'White Sneakers',
        1890.00,
        (SELECT id FROM colors WHERE name = 'White')
    )
ON CONFLICT (name) DO UPDATE
SET
    price = EXCLUDED.price,
    color_id = EXCLUDED.color_id,
    updated_at = now();

INSERT INTO clothes_categories (clothes_id, category_id)
SELECT c.id, cat.id
FROM clothes c
JOIN categories cat ON cat.name IN ('Tops', 'Casual')
WHERE c.name = 'Basic T-Shirt'
ON CONFLICT (clothes_id, category_id) DO NOTHING;

INSERT INTO clothes_categories (clothes_id, category_id)
SELECT c.id, cat.id
FROM clothes c
JOIN categories cat ON cat.name IN ('Pants', 'Casual')
WHERE c.name = 'Blue Jeans'
ON CONFLICT (clothes_id, category_id) DO NOTHING;

INSERT INTO clothes_categories (clothes_id, category_id)
SELECT c.id, cat.id
FROM clothes c
JOIN categories cat ON cat.name IN ('Outerwear', 'Casual')
WHERE c.name = 'Red Hoodie'
ON CONFLICT (clothes_id, category_id) DO NOTHING;

INSERT INTO clothes_categories (clothes_id, category_id)
SELECT c.id, cat.id
FROM clothes c
JOIN categories cat ON cat.name IN ('Shoes', 'Sport')
WHERE c.name = 'White Sneakers'
ON CONFLICT (clothes_id, category_id) DO NOTHING;
