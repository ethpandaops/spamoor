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
(1, 'eoatx', 'EOA Transaction Spammer', '800 type-2 eoa transactions per slot, gas limit 20 gwei, 16.8M gas usage', '# wallet settings
seed: eoatx-1 # seed for the wallet
refill_amount: 1000000000000000000 # refill 1 ETH when
refill_balance: 500000000000000000 # balance drops below 0.5 ETH
refill_interval: 600 # check every 10 minutes

# scenario: eoatx
total_count: 0
throughput: 800
max_pending: 1600
max_wallets: 500
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
(2, 'erctx', 'ERC20 Transaction Spammer', '430 type-2 erc20 transactions per slot, gas limit 20 gwei, 16.8M gas usage', '# wallet settings
seed: erctx-2 # seed for the wallet
refill_amount: 1000000000000000000 # refill 1 ETH when
refill_balance: 500000000000000000 # balance drops below 0.5 ETH
refill_interval: 600 # check every 10 minutes

# scenario: erctx
total_count: 0
throughput: 430
max_pending: 900
max_wallets: 500
rebroadcast: 120
base_fee: 20
tip_fee: 2
amount: 20
random_amount: false
random_target: false
', 0, 0, '{}');


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'NOT SUPPORTED';
-- +goose StatementEnd
