CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users
(
    id            BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    email         VARCHAR(320) NOT NULL UNIQUE,
    password_hash TEXT  NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp
);

CREATE TABLE accounts
(
    id         BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id    BIGINT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    balance    NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    currency   CHAR(3)        NOT NULL DEFAULT 'RUB',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT current_timestamp
);
CREATE INDEX idx_accounts_user_id ON accounts (user_id);