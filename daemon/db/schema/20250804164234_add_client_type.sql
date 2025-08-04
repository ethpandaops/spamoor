-- +goose Up
-- +goose StatementBegin

ALTER TABLE client_configs ADD COLUMN client_type TEXT NOT NULL DEFAULT 'client';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

ALTER TABLE client_configs DROP COLUMN client_type;

-- +goose StatementEnd