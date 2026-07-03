-- +goose Up
-- +goose StatementBegin

-- Add a blob member to the existing "Regular Chain Load" group so the group
-- exercises blob traffic alongside the other everyday transaction types. Uses the
-- plain "blobs" scenario (throughput-based), so it joins the group's shared
-- throughput budget by weight like the other members. Created paused (status 0)
-- like the rest of that template group. fulu_activation: 0 sends v1 (cell-proof)
-- blob sidecars by default, correct for current post-Fusaka networks.
INSERT INTO "spammers" ("id", "scenario", "name", "description", "config", "status", "created_at", "state", "group_id", "group_config")
VALUES
(30, 'blobs', 'Blob Transactions', 'Sends type-3 blob transactions as part of the shared chain load.', '# wallet settings
seed: blobs-700004 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: blobs
total_count: 0
throughput: 10
sidecars: 1
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
blob_fee: 20
blob_v1_percent: 100
fulu_activation: 0
client_group: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 9}');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM "spammers" WHERE "id" = 30 AND "group_id" = 10;

-- +goose StatementEnd
