ALTER TABLE clothes_change_logs
DROP CONSTRAINT IF EXISTS clothes_change_logs_action_check;

ALTER TABLE clothes_change_logs
ADD CONSTRAINT clothes_change_logs_action_check
CHECK (
    action IN (
        'CREATE',
        'UPDATE',
        'DELETE',
        'IMAGE_UPLOAD',
        'IMAGE_DELETE',
        'IMAGE_REPLACE',
        'IMAGE_SET_PRIMARY'
    )
);