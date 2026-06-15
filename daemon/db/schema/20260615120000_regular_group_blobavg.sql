-- +goose Up
-- +goose StatementBegin

-- Default-enable the "Regular Chain Load" group so it runs on launch. Group rows
-- are never executed directly; their members self-resume by their own status, so
-- we flip the member rows (group_id = 10) to running (status 1). The group row
-- itself (id 10) stays status 0.
UPDATE "spammers" SET "status" = 1 WHERE "group_id" = 10;

-- Add a blob-average member targeting an average of 3 blobs/block. blob-average
-- is driven by target_average (not throughput), so it joins with weight 0: it
-- takes no share of the group's shared-throughput budget and simply runs at its
-- own target average alongside the other members.
INSERT INTO "spammers" ("id", "scenario", "name", "description", "config", "status", "created_at", "state", "group_id", "group_config")
VALUES
(30, 'blob-average', 'Blob Average', 'Maintains a steady average of ~3 blobs per block.', '# wallet settings
seed: blobaverage-700004 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: blob-average
target_average: 3
sidecars: 1
max_pending: 10
max_wallets: 5
rebroadcast: 1
base_fee: 20
tip_fee: 2
blob_fee: 20
client_group: ""
log_txs: false

', 1, 0, '{}', 10, '{"weight": 0, "enabled": true, "sort_order": 9}');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM "spammers" WHERE "id" = 30 AND "group_id" = 10;
UPDATE "spammers" SET "status" = 0 WHERE "group_id" = 10;

-- +goose StatementEnd
