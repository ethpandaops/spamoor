-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "client_configs"
(
    "rpc_url" TEXT NOT NULL,
    "name" TEXT NOT NULL DEFAULT '',
    "tags" TEXT NOT NULL DEFAULT '',
    "enabled" INTEGER NOT NULL DEFAULT 1,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL,
    CONSTRAINT "client_configs_pkey" PRIMARY KEY("rpc_url")
);

CREATE INDEX IF NOT EXISTS "client_configs_enabled_idx"
    ON "client_configs"
    ("enabled" ASC);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "client_configs";

-- +goose StatementEnd