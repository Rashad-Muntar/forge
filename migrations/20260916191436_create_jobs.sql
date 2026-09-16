-- +goose Up
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    idempotency_key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'queued',
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT jobs_status_check
        CHECK (status IN (
            'queued',
            'running',
            'succeeded',
            'failed',
            'cancelled'
        )),

    CONSTRAINT jobs_priority_check
        CHECK (priority >= 0),

    CONSTRAINT jobs_attempts_check
        CHECK (attempts >= 0),

    CONSTRAINT jobs_max_attempts_check
        CHECK (max_attempts > 0),

    CONSTRAINT jobs_idempotency_key_unique
        UNIQUE (idempotency_key)
);

CREATE INDEX idx_jobs_status_available_at
    ON jobs (status, available_at);

CREATE INDEX idx_jobs_created_at
    ON jobs (created_at DESC);

-- +goose Down
DROP TABLE jobs;
