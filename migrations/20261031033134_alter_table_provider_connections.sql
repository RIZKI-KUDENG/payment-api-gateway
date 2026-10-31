-- +goose Up
SELECT 'up SQL query';
ALTER TABLE provider_connections
    ALTER COLUMN credentials_ref TYPE UUID
    USING credentials_ref::UUID;


ALTER TABLE provider_connections
    ADD CONSTRAINT fk_provider_connections_credentials
    FOREIGN KEY (tenant_id, credentials_ref)
    REFERENCES provider_credentials(tenant_id, id)
    ON DELETE RESTRICT;
-- +goose Down
SELECT 'down SQL query';

ALTER TABLE provider_connections
    DROP CONSTRAINT IF EXISTS fk_provider_connections_credentials;

ALTER TABLE provider_connections
    ALTER COLUMN credentials_ref TYPE TEXT
    USING credentials_ref::TEXT;
