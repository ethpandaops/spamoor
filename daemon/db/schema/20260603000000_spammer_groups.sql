-- +goose Up
-- +goose StatementBegin

-- Spammer groups are stored as regular "spammers" rows marked with the reserved
-- sentinel scenario name "group". Two columns link members to their parent group
-- and carry role-dependent JSON metadata:
--   group_id     - parent group row id for members (0 for standalone rows and groups)
--   group_config - JSON: for group rows {throughput_mode,total_throughput,total_count},
--                  for member rows {weight,enabled,sort_order}

ALTER TABLE "spammers" ADD COLUMN "group_id" INTEGER NOT NULL DEFAULT 0;
ALTER TABLE "spammers" ADD COLUMN "group_config" TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS "spammers_group_id_idx"
    ON "spammers"
    ("group_id" ASC);

-- The default spammer groups ("Regular Chain Load", "Fuzzing") are no longer inserted
-- via schema migrations. They are imported from the embedded YAML files in
-- daemon/default-spammers/ on first launch.

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS "spammers_group_id_idx";
ALTER TABLE "spammers" DROP COLUMN "group_config";
ALTER TABLE "spammers" DROP COLUMN "group_id";

-- +goose StatementEnd
