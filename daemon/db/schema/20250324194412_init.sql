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
random_amount: false
random_target: false
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
rebroadcast: 30
base_fee: 20
tip_fee: 2
blob_fee: 20
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
rebroadcast: 120
base_fee: 10
tip_fee: 2
gas_limit: 300000
amount: 20
data: "repeat:0x00:25000"
random_amount: false
random_target: false
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
rebroadcast: 120
base_fee: 20
tip_fee: 2
gas_units_to_burn: 2000000
client_group: ""
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
rebroadcast: 120
base_fee: 20
tip_fee: 2
pair_count: 1
min_swap_amount: "100000000000000000"
max_swap_amount: "1000000000000000000000"
buy_ratio: 50
slippage: 50
sell_threshold: "100000000000000000000000"
client_group: ""
', 0, 0, '{}');


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd
