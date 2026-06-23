UPDATE purchase_items AS pi
SET product_image_url = image_snapshot.image_url
FROM (
    SELECT DISTINCT ON (ci.clothes_id)
        ci.clothes_id,
        'http://localhost:9000/' || ci.bucket || '/' || ci.object_key AS image_url
    FROM clothes_images AS ci
    ORDER BY ci.clothes_id, ci.id ASC
) AS image_snapshot
WHERE pi.clothes_id = image_snapshot.clothes_id
  AND pi.product_image_url IS NULL;
