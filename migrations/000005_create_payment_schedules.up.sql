CREATE TABLE payment_schedules
(
    id         BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    credit_id  BIGINT         NOT NULL REFERENCES credits (id) ON DELETE CASCADE,
    due_date   DATE           NOT NULL,
    amount     NUMERIC(12, 2) NOT NULL,
    paid       BOOLEAN        NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_schedule_credit_id ON payment_schedules (credit_id);
