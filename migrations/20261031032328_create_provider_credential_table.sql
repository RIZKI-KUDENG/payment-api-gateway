-- +goose Up

CREATE TABLE IF NOT EXISTS provider_credentials (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    tenant_id UUID NOT NULL
        REFERENCES tenants(id) ON DELETE CASCADE,

    encrypted_value BYTEA NOT NULL,

    encryption_key_version SMALLINT NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_provider_credential_tenant
        UNIQUE (tenant_id, id)
);

-- +goose Down

DROP TABLE IF EXISTS provider_credentials;
