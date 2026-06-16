-- +goose Up
-- +goose StatementBegin

-- "Fuzzing" group: bundles the EVM execution fuzzer (evm-fuzz), the transaction-
-- layer fuzzer (tx-fuzz) and the invalid-transaction fuzzer (tx-fuzz-invalid)
-- under one shared throughput budget. Created paused (status 0) as a template,
-- mirroring the "Regular Chain Load" group. Reserved low ids (20-23) never
-- collide with runtime-created spammers/groups, which start at id 100.
INSERT INTO "spammers" ("id", "scenario", "name", "description", "config", "status", "created_at", "state", "group_id", "group_config")
VALUES
(20, 'group', 'Fuzzing', 'A mix of fuzzing scenarios - EVM execution fuzzing (random bytecode deployments) and transaction-layer fuzzing (randomized tx types, calldata, access lists, authorizations and blobs) - sharing a global budget of 10 tx/slot to surface consensus and validation bugs.', '', 0, 0, '{}', 0, '{"throughput_mode": "shared", "total_throughput": 10, "total_count": 0, "total_max_pending": 0}'),
(21, 'evm-fuzz', 'EVM Execution Fuzzing', 'Deploys contracts with randomly generated, stack-aware bytecode exercising opcodes and precompiles.', '# wallet settings
seed: evmfuzz-700001 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: evm-fuzz
total_count: 0
throughput: 4
max_pending: 20
max_wallets: 4
rebroadcast: 1
base_fee: 20
tip_fee: 2
gas_limit: 1000000
max_code_size: 512
min_code_size: 100
fuzz_mode: all
client_group: ""
log_txs: false

', 0, 0, '{}', 20, '{"weight": 2, "enabled": true, "sort_order": 0}'),
(22, 'tx-fuzz', 'Transaction Layer Fuzzing', 'Sends well-formed transactions across all types (legacy/2930/1559/4844/7702) with randomized calldata, access lists, authorizations, blobs and targets.', '# wallet settings
seed: txfuzz-700002 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: tx-fuzz
total_count: 0
throughput: 4
max_pending: 20
max_wallets: 4
rebroadcast: 1
base_fee: 20
tip_fee: 2
gas_limit: 500000
tx_types: all
max_call_data: 1024
max_access_list: 5
max_auth_list: 5
max_blobs: 3
client_group: ""
log_txs: false

', 0, 0, '{}', 20, '{"weight": 2, "enabled": true, "sort_order": 1}'),
(23, 'tx-fuzz-invalid', 'Invalid Transaction Fuzzing', 'Fires deliberately-invalid transactions (bad chainid/nonce/gas/fees, malformed RLP, empty 7702 auth, blobless blob tx) from a small reused pool of burner wallets, out-of-pool and fire-and-forget. A node accepting one is a potential validation gap.', '# wallet settings
seed: txfuzzinvalid-700003 # seed for the wallet
refill_amount: 50000000000000000 # refill 0.05 ETH when
refill_balance: 10000000000000000 # balance drops below 0.01 ETH
refill_interval: 600 # check every 10 minutes

# scenario: tx-fuzz-invalid
total_count: 0
throughput: 2
max_pending: 20
max_wallets: 1
categories: all
client_group: ""
log_txs: false

', 0, 0, '{}', 20, '{"weight": 1, "enabled": true, "sort_order": 2}');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM "spammers" WHERE "id" IN (21, 22, 23) AND "group_id" = 20;
DELETE FROM "spammers" WHERE "id" = 20 AND "scenario" = 'group';

-- +goose StatementEnd
