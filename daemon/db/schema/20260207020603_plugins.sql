-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "plugins"
(
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL UNIQUE,
    "source_type" TEXT NOT NULL,
    "source_path" TEXT NOT NULL,
    "metadata_name" TEXT,
    "metadata_build_time" TEXT,
    "metadata_git_version" TEXT,
    "archive_data" TEXT,
    "scenarios" TEXT,
    "enabled" INTEGER NOT NULL DEFAULT 1,
    "load_error" TEXT,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS "plugins_name_idx"
    ON "plugins"
    ("name" ASC);

CREATE INDEX IF NOT EXISTS "plugins_enabled_idx"
    ON "plugins"
    ("enabled" ASC);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS "plugins_enabled_idx";
DROP INDEX IF EXISTS "plugins_name_idx";
DROP TABLE IF EXISTS "plugins";

-- +goose StatementEnd
