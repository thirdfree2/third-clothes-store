CREATE TABLE clothes_change_logs (
    id BIGSERIAL PRIMARY KEY,

    clothes_id BIGINT NULL REFERENCES clothes(id) ON DELETE SET NULL,

    action VARCHAR(20) NOT NULL,
    before_data JSONB NULL,
    after_data JSONB NULL,

    changed_by VARCHAR(100) NULL,
    note TEXT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT clothes_change_logs_action_check
        CHECK (action IN ('CREATE', 'UPDATE', 'DELETE'))
);

CREATE INDEX idx_clothes_change_logs_clothes_id
    ON clothes_change_logs(clothes_id);

CREATE INDEX idx_clothes_change_logs_created_at
    ON clothes_change_logs(created_at);