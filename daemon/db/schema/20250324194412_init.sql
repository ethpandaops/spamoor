-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "spamoor_state"
(
    "key" TEXT NOT NULL UNIQUE,
    "value" TEXT,
    CONSTRAINT "spamoor_state_pkey" PRIMARY KEY("key")
);

CREATE TABLE IF NOT EXISTS "spammers"
(
    "id" INTEGER NOT NULL,
    "scenario" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "config" TEXT NOT NULL,
    "status" INTEGER NOT NULL DEFAULT 0,
    "created_at" INTEGER NOT NULL,
    "state" TEXT NOT NULL,
    CONSTRAINT "spammers_pkey" PRIMARY KEY("id")
);

CREATE INDEX IF NOT EXISTS "spammers_scenario_idx"
    ON "spammers"
    ("scenario" ASC);

INSERT INTO "spammers" ("id", "scenario", "name", "description", "config", "status", "created_at", "state")
VALUES 
-- EOA Transaction Spammer
(1, 'eoatx', 'EOA Transaction Spammer', '400 type-2 eoa transactions per slot, gas limit 20 gwei, 8.4M gas usage', '# wallet settings
seed: eoatx-1 # seed for the wallet
refill_amount: 1000000000000000000 # refill 1 ETH when
refill_balance: 500000000000000000 # balance drops below 0.5 ETH
refill_interval: 600 # check every 10 minutes

# scenario: eoatx
total_count: 0
throughput: 400
max_pending: 800
max_wallets: 400
rebroadcast: 120
base_fee: 20
tip_fee: 2
gas_limit: 21000
amount: 20
data: ""
to: ""
timeout: ""
random_amount: false
random_target: false
self_tx_only: false
client_group: ""
log_txs: false
', 0, 0, '{}'),

-- ERC20 Transaction Spammer
(2, 'erctx', 'ERC20 Transaction Spammer', '215 type-2 erc20 transactions per slot, gas limit 20 gwei, 8.4M gas usage', '# wallet settings
seed: erctx-2 # seed for the wallet
refill_amount: 1000000000000000000 # refill 1 ETH when
refill_balance: 500000000000000000 # balance drops below 0.5 ETH
refill_interval: 600 # check every 10 minutes

# scenario: erctx
total_count: 0
throughput: 215
max_pending: 430
max_wallets: 400
rebroadcast: 120
base_fee: 20
tip_fee: 2
amount: 20
random_amount: false
random_target: false
timeout: ""
client_group: ""
deploy_client_group: ""
log_txs: false
', 0, 0, '{}'),

-- Blob Transaction Spammer
(3, 'blob-combined', 'Blob Transaction Spammer', '3 type-4 blob transactions per block with 1-2 blobs', '# wallet settings
seed: blob-combined-3 # seed for the wallet
refill_amount: 2000000000000000000 # refill 2 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: blob-combined
total_count: 0
throughput: 3
sidecars: 2
max_pending: 6
max_wallets: 10
replace: 30
max_replacements: 4
rebroadcast: 1
base_fee: 20
tip_fee: 2
blob_fee: 20
blob_v1_percent: 100
fulu_activation: 9223372036854775807
throughput_increment_interval: 0
timeout: ""
client_group: ""
log_txs: false
', 0, 0, '{}'),

-- Big Block Spammer
(4, 'eoatx', 'Big Block Spammer', '200 type-2 eoa transactions per slot with 25k zero bytes calldata each', '# wallet settings
seed: eoatx-4 # seed for the wallet
refill_amount: 2000000000000000000 # refill 2 ETH when
refill_balance: 500000000000000000 # balance drops below 0.5 ETH
refill_interval: 600 # check every 10 minutes

# scenario: eoatx
total_count: 0
throughput: 200
max_pending: 400
max_wallets: 200
rebroadcast: 1
base_fee: 10
tip_fee: 2
gas_limit: 50000
amount: 20
data: ""
to: ""
timeout: ""
random_amount: false
random_target: false
self_tx_only: false
client_group: ""
log_txs: false
', 0, 0, '{}'),

-- Gas Burner Spammer
(5, 'gasburnertx', 'Gas Burner Spammer', '20 gas-burner transactions per slot, burning 2M gas each, gas limit 20 gwei, up to 40M gas usage', '# wallet settings
seed: gasburnertx-5 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: gasburnertx
total_count: 0
throughput: 20
max_pending: 20
max_wallets: 40
rebroadcast: 1
base_fee: 20
tip_fee: 2
gas_units_to_burn: 2000000
gas_remainder: 10000
timeout: ""
opcodes: ""
init_opcodes: ""
client_group: ""
deploy_client_group: ""
log_txs: false
', 0, 0, '{}'),

-- Uniswap V2 Swap Spammer
(6, 'uniswap-swaps', 'Uniswap V2 Swap Spammer', '200 type-2 swap transactions with uniswap v2 pools, gas limit 20 gwei, ~15M gas usage', '# wallet settings
seed: uniswap-swaps-6 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: uniswap-swaps
total_count: 0
throughput: 200
max_pending: 200
max_wallets: 200
rebroadcast: 1
base_fee: 20
tip_fee: 2
pair_count: 1
min_swap_amount: "100000000000000000"
max_swap_amount: "1000000000000000000000"
buy_ratio: 40
slippage: 50
sell_threshold: "50000000000000000000000"
timeout: ""
client_group: ""
deploy_client_group: ""
log_txs: false
', 0, 0, '{}'),

-- EVM Fuzz Spammer
(7, 'evm-fuzz', 'EVM Fuzz Spammer', 'Opcode & Precompile fuzzer, 50tx/block, gas limit 1M, ~20M gas usage', '# wallet settings
seed: evm-fuzz-7 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: evm-fuzz
total_count: 0
throughput: 50
max_pending: 100
max_wallets: 0
rebroadcast: 30
base_fee: 20
tip_fee: 2
gas_limit: 1000000
timeout: ""
client_group: ""
log_txs: false
max_code_size: 512
min_code_size: 100
payload_seed: ""
tx_id_offset: 0
fuzz_mode: all


', 0, 0, '{}');


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd
