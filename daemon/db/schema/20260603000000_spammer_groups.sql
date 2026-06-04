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

-- Default "Regular Chain Load" group: a balanced mix of common transaction types sharing
-- a global throughput budget of 100 tx/slot. Created paused (status 0) as a template.
INSERT INTO "spammers" ("id", "scenario", "name", "description", "config", "status", "created_at", "state", "group_id", "group_config")
VALUES
(10, 'group', 'Regular Chain Load', 'A balanced mix of everyday transaction types - EOA transfers, ERC-20/721/1155 tokens, ERC-4337 account abstraction, Safe multisig, and Curve / Uniswap V2+V3 swaps - sharing a global budget of 100 tx/slot to emulate organic chain activity.', '', 0, 0, '{}', 0, '{"throughput_mode": "shared", "total_throughput": 100, "total_count": 0, "total_max_pending": 0}'),
(11, 'eoatx', 'EOA Transfers', 'Plain ETH value transfers between externally-owned accounts.', '# wallet settings
seed: eoatx-131922 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: eoatx
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
gas_limit: 0
amount: 200
data: ""
to: ""
timeout: ""
random_amount: true
random_target: false
self_tx_only: false
client_group: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 2, "enabled": true, "sort_order": 0}'),
(12, 'erc20tx', 'ERC-20 Token Transfers', 'Deploys an ERC-20 token and continuously transfers it between wallets.', '# wallet settings
seed: erc20tx-124868 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: erc20tx
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
amount: 1000
random_amount: true
random_target: false
timeout: ""
client_group: ""
deploy_client_group: ""
token_seed: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 1}'),
(13, 'erc721tx', 'ERC-721 NFT Activity', 'Mints and transfers ERC-721 non-fungible tokens.', '# wallet settings
seed: erc721tx-849573 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: erc721tx
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
max_index: 0
random_index: false
random_target: false
timeout: ""
client_group: ""
deploy_client_group: ""
token_seed: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 2}'),
(14, 'erc1155tx', 'ERC-1155 Multi-Token Activity', 'Mints and transfers ERC-1155 multi-tokens.', '# wallet settings
seed: erc1155tx-557717 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: erc1155tx
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
amount: 100
max_index: 0
batch_size: 2
random_amount: true
random_target: false
random_index: false
random_batch_size: true
timeout: ""
client_group: ""
deploy_client_group: ""
token_seed: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 3}'),
(15, 'erc4337', 'ERC-4337 Account Abstraction', 'Submits ERC-4337 UserOperations through the EntryPoint via a bundler.', '# wallet settings
seed: erc4337-426060 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: erc4337
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
bundle_size: 1
new_account_interval: 1000
paymaster_deposit: 10
timeout: ""
client_group: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 4}'),
(16, 'safe-multisig', 'Safe Multisig Executions', 'Executes transactions through Safe (Gnosis) multisig wallets.', '# wallet settings
seed: safe-multisig-96581 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: safe-multisig
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
min_owners: 1
max_owners: 5
threshold: 0
safes_per_wallet: 3
contract_ratio: 0.5
recreate_rate: 0.05
burn_rounds: 1000
eoa_value: 1000
funding_interval: 64
gas_limit: 0
timeout: ""
client_group: ""
deploy_client_group: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 5}'),
(17, 'curve-swaps', 'Curve StableSwap Swaps', 'Performs token swaps on Curve StableSwap pools.', '# wallet settings
seed: curve-swaps-222045 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: curve-swaps
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
pool_count: 1
amplification: 200
fee: 4000000
seed_amount: "1000000000000000000000000"
wallet_funding: "10000000000000000000000"
min_swap_amount: "1000000000000000000"
max_swap_amount: "1000000000000000000000"
slippage: 100
timeout: ""
client_group: ""
deploy_client_group: ""
log_txs: false

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 6}'),
(18, 'uniswap-swaps', 'Uniswap V2 Swaps', 'Swaps tokens across Uniswap V2 pairs.', '# wallet settings
seed: uniswap-swaps-803935 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: uniswap-swaps
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
version: 2 # Uniswap V2
fee_tier: 3000
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

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 7}'),
(19, 'uniswap-swaps', 'Uniswap V3 Swaps', 'Swaps tokens across Uniswap V3 pools (0.3% fee tier).', '# wallet settings
seed: uniswap-swaps-745817 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes
# fee_strategy: adaptive # uncomment for dynamic fee headroom with normal distribution

# scenario: uniswap-swaps
total_count: 0
throughput: 10
max_pending: 20
max_wallets: 10
rebroadcast: 1
base_fee: 20
tip_fee: 2
base_fee_wei: ""
tip_fee_wei: ""
version: 3 # Uniswap V3
fee_tier: 3000
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

', 0, 0, '{}', 10, '{"weight": 1, "enabled": true, "sort_order": 8}');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DELETE FROM "spammers" WHERE "id" BETWEEN 10 AND 19;
DROP INDEX IF EXISTS "spammers_group_id_idx";
ALTER TABLE "spammers" DROP COLUMN "group_config";
ALTER TABLE "spammers" DROP COLUMN "group_id";

-- +goose StatementEnd
