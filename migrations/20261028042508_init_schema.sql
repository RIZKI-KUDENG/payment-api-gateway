-- +goose Up

CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(80) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS merchant_accounts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    name VARCHAR(80) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_merchant_tenant_id UNIQUE (tenant_id, id)
);
CREATE INDEX IF NOT EXISTS idx_merchant_accounts_tenant_id ON merchant_accounts(tenant_id);

CREATE TABLE IF NOT EXISTS provider_connections (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    provider VARCHAR(30) NOT NULL,
    provider_account_id VARCHAR(150) NOT NULL,
    name VARCHAR(80) NOT NULL,
    credentials_ref TEXT NOT NULL,
    environment VARCHAR(20) NOT NULL DEFAULT 'sandbox',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_provider_connection_tenant_id UNIQUE (tenant_id, id)
);
CREATE INDEX IF NOT EXISTS idx_provider_connections_tenant_id ON provider_connections(tenant_id);

CREATE TABLE IF NOT EXISTS merchant_provider_connections (
    tenant_id UUID NOT NULL,
    merchant_account_id UUID NOT NULL,
    provider_connection_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, merchant_account_id, provider_connection_id),
    CONSTRAINT fk_mpc_merchant FOREIGN KEY (tenant_id, merchant_account_id)
        REFERENCES merchant_accounts (tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT fk_mpc_provider FOREIGN KEY (tenant_id, provider_connection_id)
        REFERENCES provider_connections (tenant_id, id) ON DELETE RESTRICT
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_mpc_default_per_merchant
    ON merchant_provider_connections (tenant_id, merchant_account_id)
    WHERE is_default = TRUE;

CREATE TABLE IF NOT EXISTS payment_intents (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    merchant_account_id UUID NOT NULL,
    merchant_reference VARCHAR(150) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    status VARCHAR(20) NOT NULL,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_payment_intent_tenant_id UNIQUE (tenant_id, id),
    CONSTRAINT fk_payment_merchant FOREIGN KEY (tenant_id, merchant_account_id)
        REFERENCES merchant_accounts (tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT uq_payment_merchant_reference UNIQUE (tenant_id, merchant_account_id, merchant_reference)
);
CREATE INDEX IF NOT EXISTS idx_payment_intents_tenant_merchant ON payment_intents(tenant_id, merchant_account_id);

CREATE TABLE IF NOT EXISTS payment_attempts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    payment_intent_id UUID NOT NULL,
    provider_connection_id UUID NOT NULL,
    attempt_number INT NOT NULL CHECK (attempt_number > 0),
    provider_payment_id VARCHAR(150),
    provider_status VARCHAR(20),
    status VARCHAR(20) NOT NULL,
    error_code VARCHAR(100),
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payment_attempts_intent FOREIGN KEY (tenant_id, payment_intent_id)
        REFERENCES payment_intents (tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT fk_payment_attempts_provider FOREIGN KEY (tenant_id, provider_connection_id)
        REFERENCES provider_connections (tenant_id, id) ON DELETE RESTRICT,
    CONSTRAINT uq_payment_attempt_number UNIQUE (payment_intent_id, attempt_number)
);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_intent_id ON payment_attempts(payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_provider ON payment_attempts(provider_connection_id, provider_payment_id);

CREATE TABLE IF NOT EXISTS webhook_events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    provider_connection_id UUID NOT NULL REFERENCES provider_connections(id) ON DELETE RESTRICT,
    provider_event_id VARCHAR(200) NOT NULL,
    event_fingerprint VARCHAR(128) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload_ref TEXT,
    processing_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    received_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT,
    CONSTRAINT uq_webhook_provider_event UNIQUE (provider_connection_id, provider_event_id)
);
CREATE INDEX IF NOT EXISTS idx_webhook_events_fingerprint ON webhook_events(event_fingerprint);

CREATE TABLE IF NOT EXISTS idempotency_records (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    idempotency_key VARCHAR(255) NOT NULL,
    request_hash VARCHAR(128) NOT NULL,
    response_ref TEXT,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_idempotency_tenant_key UNIQUE (tenant_id, idempotency_key)
);

CREATE TABLE IF NOT EXISTS settlement_batches (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    provider_connection_id UUID NOT NULL REFERENCES provider_connections(id) ON DELETE RESTRICT,
    provider_batch_id VARCHAR(200) NOT NULL,
    settlement_date TIMESTAMPTZ NOT NULL,
    currency CHAR(3) NOT NULL,
    gross_amount BIGINT NOT NULL,
    fee BIGINT NOT NULL,
    net_amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_settlement_provider_batch UNIQUE (provider_connection_id, provider_batch_id)
);

CREATE TABLE IF NOT EXISTS settlement_items (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    settlement_batch_id UUID NOT NULL REFERENCES settlement_batches(id) ON DELETE RESTRICT,
    provider_transaction_id VARCHAR(200),
    provider_payment_id VARCHAR(150),
    matched_payment_id UUID REFERENCES payment_intents(id) ON DELETE SET NULL,
    gross_amount BIGINT NOT NULL,
    fee BIGINT NOT NULL,
    net_amount BIGINT NOT NULL,
    currency CHAR(3) NOT NULL,
    match_status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_settlement_items_batch ON settlement_items(settlement_batch_id);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    actor_id VARCHAR(150),
    actor_type VARCHAR(30) NOT NULL,
    source VARCHAR(30) NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(150),
    request_id VARCHAR(150),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id ON audit_logs(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS settlement_items;
DROP TABLE IF EXISTS settlement_batches;
DROP TABLE IF EXISTS idempotency_records;
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS payment_attempts;
DROP TABLE IF EXISTS payment_intents;
DROP TABLE IF EXISTS merchant_provider_connections;
DROP TABLE IF EXISTS provider_connections;
DROP TABLE IF EXISTS merchant_accounts;
DROP TABLE IF EXISTS tenants;
