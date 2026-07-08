# Spamoor API Consumer Guide

This guide provides comprehensive documentation for consuming the Spamoor daemon API programmatically. The API allows you to create, manage, and monitor transaction spammers remotely.

## Table of Contents

- [API Overview](#api-overview)
- [Authentication](#authentication)
- [Spammer Management](#spammer-management)
- [Scenario Management](#scenario-management)
- [Client Management](#client-management)
- [Real-time Monitoring](#real-time-monitoring)
- [Import/Export](#importexport)
- [System Endpoints](#system-endpoints)
- [Error Handling](#error-handling)
- [Example Workflows](#example-workflows)
- [SDKs and Libraries](#sdks-and-libraries)

## API Overview

The Spamoor daemon exposes a RESTful API on port 8080 (configurable) with the following characteristics:

- **Base URL**: `http://localhost:8080/api`
- **Content Type**: `application/json` (for most endpoints)
- **Authentication**: Optional bearer-token JWT verified against an `authenticatoor` service (see below). Disabled by default.
- **Documentation**: Interactive Swagger UI at `/docs`
- **Metrics**: Prometheus metrics at `/metrics`

## Authentication

Authentication is opt-in. When the daemon is started without `--auth-provider-url`, the API is open and every request is accepted — suitable only for trusted networks.

When `--auth-provider-url=<url>` is set, all write endpoints (and SSE streams) require a JWT in the `Authorization: Bearer <token>` header. Tokens are verified against the JWKS published by the configured [`service-authenticatoor`](https://github.com/ethpandaops/service-authenticatoor) instance:

- **Issuer (`iss`)**: discovered from `<url>/.well-known/openid-configuration` (falls back to `<url>` itself).
- **Audience (`aud`)**: the parent DNS zone of the auth service host (e.g. `auth.foo.example.io` → `foo.example.io`).
- **Identity for audit logs**: the token's `email` claim, falling back to `sub`.

For SSE endpoints that cannot set headers, the token may be passed as the `auth` query parameter (sent as `Authorization: Bearer <token>` server-side).

```bash
# Open mode (default): no header required
curl http://localhost:8080/api/spammers

# With auth-provider-url configured
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/spammers
```

Browser users are redirected to the authenticatoor's login page by the web UI; the token is then stored client-side and attached automatically.

## Spammer Management

### List All Spammers

Get a list of all configured spammers.

```http
GET /api/spammers
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "EOA Test Spammer",
    "description": "Testing basic EOA transactions",
    "scenario": "eoatx",
    "status": 1,
    "created_at": "2024-01-15T10:30:00.123456789Z"
  }
]
```

**Status Values:**
- `0`: Stopped
- `1`: Running
- `2`: Paused

### Create a New Spammer

Create a new spammer with specified configuration.

```http
POST /api/spammer
Content-Type: application/json

{
  "name": "My EOA Spammer",
  "description": "High-throughput EOA transactions",
  "scenario": "eoatx",
  "config": "throughput: 10\nmax_pending: 50\namount: 100\nrandom_amount: true",
  "startImmediately": true
}
```

**Response:**
```json
2
```
Returns the new spammer's ID.

### Get Spammer Details

Get detailed information about a specific spammer.

```http
GET /api/spammer/{id}
```

**Response:**
```json
{
  "id": 1,
  "name": "EOA Test Spammer",
  "description": "Testing basic EOA transactions",
  "scenario": "eoatx",
  "config": "throughput: 10\nmax_pending: 50\namount: 100",
  "status": 1
}
```

### Update Spammer Configuration

Update an existing spammer's configuration.

```http
PUT /api/spammer/{id}
Content-Type: application/json

{
  "name": "Updated EOA Spammer",
  "description": "Modified configuration",
  "config": "throughput: 15\nmax_pending: 100\namount: 200"
}
```

### Start a Spammer

Start a paused or stopped spammer.

```http
POST /api/spammer/{id}/start
```

### Pause a Spammer

Pause a running spammer.

```http
POST /api/spammer/{id}/pause
```

### Delete a Spammer

Delete a spammer (stops it if running).

```http
DELETE /api/spammer/{id}
```

### Reclaim Funds

Reclaim funds from a spammer's wallet pool back to the root wallet.

```http
POST /api/spammer/{id}/reclaim
```

## Spammer Groups

A **spammer group** lets you control many spammers as one unit. A group applies a
shared **overlay** of common config fields on top of each member's own full config, and
can optionally split a **global throughput** (or total count) budget across members by
weight.

A group is stored as a regular spammer row with the reserved sentinel `scenario: "group"`.
Members keep their own full config; adding/removing a member only sets/clears its
`group_id`, so an ungrouped member behaves exactly like a standalone spammer.

List entries (`GET /api/spammers`) and `GET /api/spammer/{id}` include `is_group`,
`group_id`, and a parsed `group_config` (for groups) or `member_config` (for members).
Group rows report a **derived status** (running if any enabled member runs).

The lifecycle endpoints (`/start`, `/pause`, `/reclaim`) work on a group id and fan out
to its enabled members.

### Create a Group

```http
POST /api/spammer-group
Content-Type: application/json

{
  "name": "mixed-load",
  "description": "broad transaction variance",
  "config": "base_fee: 20\ntip_fee: 2\n",
  "throughput_mode": "shared",
  "total_throughput": 100,
  "total_count": 0,
  "total_max_pending": 0,
  "auto_restart_failed": false,
  "auto_restart_cooldown": 0
}
```

`throughput_mode` is `"shared"` (split totals by weight) or `"independent"` (overlay
only). In shared mode `total_throughput`/`total_count` are split across enabled members;
each member's `max_wallets` is derived from its resolved throughput (≈ throughput/4,
clamped to [20, 1000]) and its `max_pending` defaults to 2× its resolved throughput.
Set `total_max_pending` to a non-zero value to instead split an explicit
concurrent-pending budget across members by weight. Returns the new group id.

With `auto_restart_failed` enabled, members that stop in the failed (error) state are
restarted automatically after `auto_restart_cooldown` seconds (0 = default 300).
Members stopped normally (paused or finished) are never restarted, and manual actions
taken during the cooldown (restart, pause, removal) always win.

### Update a Group

```http
PUT /api/spammer-group/{id}
Content-Type: application/json

{ "name": "...", "description": "...", "config": "...", "throughput_mode": "shared", "total_throughput": 200, "total_count": 0 }
```

### Add / Update a Member

Assign a spammer to a group (or update its weight/enabled/order). A scenario with no
`throughput`/`total_count` field is rejected from a shared group.

```http
PUT /api/spammer/{id}/group
Content-Type: application/json

{ "group_id": 100, "weight": 20, "enabled": true, "sort_order": 0 }
```

### Detach a Member

Remove a spammer from its group, leaving it a working standalone spammer.

```http
DELETE /api/spammer/{id}/group
```

### Reorder Members

```http
PUT /api/spammer-group/{id}/members/order
Content-Type: application/json

{ "order": [101, 103, 102] }
```

### Delete a Group

For a group, the `cascade` query parameter controls member handling: `cascade=false`
(default) detaches members to standalone spammers; `cascade=true` deletes them too.

```http
DELETE /api/spammer/{id}?cascade=false
```

### Preview a Member's Effective Config

Returns the config a member will actually run with (group overlay + resolved
throughput/count split + derived `max_wallets`).

```http
GET /api/spammer/{id}/effective-config
```

## Scenario Management

### List Available Scenarios

Get all available transaction scenarios.

```http
GET /api/scenarios
```

**Response:**
```json
[
  {
    "name": "eoatx",
    "description": "Send standard EOA transactions with configurable amounts and targets"
  },
  {
    "name": "blobs",
    "description": "Send blob transactions with random data"
  }
]
```

### Get Scenario Configuration Template

Get a default YAML configuration template for a specific scenario.

```http
GET /api/scenarios/{name}/config
```

**Response:**
```yaml
# wallet settings
seed: eoatx-123456 # seed for the wallet
refill_amount: 5000000000000000000 # refill 5 ETH when
refill_balance: 1000000000000000000 # balance drops below 1 ETH
refill_interval: 600 # check every 10 minutes

# scenario: eoatx
throughput: 0
count: 0
max_pending: 0
# ... scenario-specific options
```

## Client Management

### List All Clients

Get information about all RPC clients.

```http
GET /api/clients
```

**Response:**
```json
[
  {
    "index": 0,
    "name": "geth-1",
    "group": "mainnet",
    "groups": ["mainnet", "primary"],
    "version": "Geth/v1.13.8-stable",
    "block_height": 19234567,
    "ready": true,
    "rpc_host": "http://localhost:8545",
    "enabled": true,
    "name_override": "Custom Node Name"
  }
]
```

### Update Client Groups

Update the group assignment for a client.

```http
PUT /api/client/{index}/group
Content-Type: application/json

{
  "groups": ["mainnet", "backup", "testing"]
}
```

For backward compatibility, single group is also supported:
```json
{
  "group": "mainnet"
}
```

### Enable/Disable Client

Control whether a client is used for transactions.

```http
PUT /api/client/{index}/enabled
Content-Type: application/json

{
  "enabled": false
}
```

### Update Client Name

Set a custom display name for a client.

```http
PUT /api/client/{index}/name
Content-Type: application/json

{
  "name_override": "Primary Mainnet Node"
}
```

## Real-time Monitoring

### Get Spammer Logs

Retrieve recent log entries for a spammer.

```http
GET /api/spammer/{id}/logs
```

**Response:**
```json
[
  {
    "time": "2024-01-15T10:30:15.123456789Z",
    "level": "info",
    "message": "transaction submitted",
    "fields": {
      "hash": "0x1234567890abcdef...",
      "nonce": "42",
      "wallet": "0xabcdef..."
    }
  }
]
```

### Stream Real-time Logs

Use Server-Sent Events (SSE) to stream logs in real-time.

```http
GET /api/spammer/{id}/logs/stream?since=2024-01-15T10:30:00.000000000Z
Accept: text/event-stream
```

**JavaScript Example:**
```javascript
const eventSource = new EventSource('/api/spammer/1/logs/stream');

eventSource.onmessage = function(event) {
  const logs = JSON.parse(event.data);
  logs.forEach(log => {
    console.log(`[${log.level}] ${log.message}`, log.fields);
  });
};

eventSource.onerror = function(event) {
  console.error('SSE connection error:', event);
};
```

**Bash/curl Example:**
```bash
# Stream logs in real-time
curl -N "http://localhost:8080/api/spammer/1/logs/stream?since=2024-01-15T10:30:00.000000000Z" \
  -H "Accept: text/event-stream" | while IFS= read -r line; do
    if [[ $line == data:* ]]; then
        # Extract JSON data after "data: " prefix
        json_data="${line#data: }"
        echo "$json_data" | jq -r '.[] | "[\(.level)] \(.message)"'
    fi
done
```

## Import/Export

### Export Spammers

Export spammer configurations to YAML format.

```http
POST /api/spammers/export
Content-Type: application/json

{
  "spammer_ids": [1, 2, 3]
}
```

To export all spammers, send an empty array or omit the field:
```json
{
  "spammer_ids": []
}
```

**Response:**
```yaml
- scenario: eoatx
  name: "EOA Test Spammer"
  description: "Testing basic EOA transactions"
  config:
    throughput: 10
    max_pending: 50
    amount: 100
```

Groups are exported as a `scenario: group` entry (its `config` is the sparse overlay and
`group_config` carries the mode/totals), and each member carries a `group` back-reference
plus its `group_config` (weight/enabled/sort_order). The parent group of any exported
member is always included so the import can re-link it:

```yaml
- scenario: group
  name: mixed-load
  config:
    base_fee: 20
  group_config:
    throughput_mode: shared
    total_throughput: 100
    total_count: 0
- scenario: eoatx
  name: eoa-member
  group: mixed-load
  config:
    amount: 100
  group_config:
    weight: 20
    enabled: true
    sort_order: 0
```

On import, groups are created first, then members are linked by group name. Importing the
same YAML through the CLI `spamoor run` path resolves groups in-process (overlay +
weight split) so members behave identically to the daemon.

### Import Spammers

Import spammers from YAML data or URL.

```http
POST /api/spammers/import
Content-Type: application/json

{
  "input": "- scenario: eoatx\n  name: Imported Spammer\n  config:\n    throughput: 5"
}
```

Or import from URL:
```json
{
  "input": "https://example.com/spammer-config.yaml"
}
```

**Response:**
```json
{
  "data": {
    "imported": 2,
    "skipped": 0,
    "errors": []
  }
}
```

## System Endpoints

### Prometheus Metrics

Get system metrics for monitoring.

```http
GET /metrics
```

**Response:**
```
# HELP spamoor_spammer_running Number of running spammers
# TYPE spamoor_spammer_running gauge
spamoor_spammer_running 3

# HELP spamoor_transactions_sent_total Total number of transactions sent
# TYPE spamoor_transactions_sent_total counter
spamoor_transactions_sent_total{scenario="eoatx",spammer="1"} 1234

# HELP spamoor_transaction_failures_total Total number of failed transactions
# TYPE spamoor_transaction_failures_total counter
spamoor_transaction_failures_total{scenario="eoatx",spammer="1"} 5
```

### API Documentation

Interactive Swagger UI documentation.

```http
GET /docs/
```

## Error Handling

The API uses standard HTTP status codes and returns JSON error responses:

### Common Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

### Error Response Format

```json
{
  "error": "Spammer not found"
}
```

### Example Error Handling

```bash
# Function to handle API responses with error checking
call_api() {
    local method="$1"
    local url="$2"
    local data="$3"
    
    if [[ -n "$data" ]]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [[ $http_code -ge 400 ]]; then
        error_msg=$(echo "$body" | jq -r '.error // "Unknown error"')
        echo "Error: HTTP $http_code - $error_msg" >&2
        return 1
    fi
    
    echo "$body"
}

# Usage example
spammer_config='{
  "name": "Test Spammer",
  "scenario": "eoatx",
  "config": "throughput: 5",
  "startImmediately": false
}'

if spammer_id=$(call_api "POST" "http://localhost:8080/api/spammer" "$spammer_config"); then
    echo "Created spammer with ID: $spammer_id"
else
    echo "Failed to create spammer"
fi
```

## Example Workflows

### Creating and Managing a Spammer

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api"

# Source the error handling function
source error_handling.sh  # Contains the call_api function from above

# 1. List available scenarios
echo "=== Available Scenarios ==="
scenarios=$(call_api "GET" "$BASE_URL/scenarios")
echo "$scenarios" | jq -r '.[] | "- \(.name): \(.description)"'
echo

# 2. Get configuration template
echo "=== EOA Transaction Template ==="
default_config=$(curl -s "$BASE_URL/scenarios/eoatx/config")
echo "$default_config"
echo

# 3. Create a spammer
echo "=== Creating Spammer ==="
spammer_config='{
  "name": "Test EOA Spammer",
  "description": "API-created spammer for testing",
  "scenario": "eoatx",
  "config": "throughput: 5\nmax_pending: 20\namount: 100\nrandom_amount: true",
  "startImmediately": false
}'

spammer_id=$(call_api "POST" "$BASE_URL/spammer" "$spammer_config")
spammer_id=$(echo "$spammer_id" | tr -d '"')  # Remove quotes
echo "Created spammer with ID: $spammer_id"

# 4. Start the spammer
echo "=== Starting Spammer ==="
call_api "POST" "$BASE_URL/spammer/$spammer_id/start" >/dev/null
echo "Spammer started"

# 5. Monitor for 30 seconds
echo "=== Monitoring Logs for 30 seconds ==="
timeout 30s curl -N "$BASE_URL/spammer/$spammer_id/logs/stream" \
  -H "Accept: text/event-stream" | while IFS= read -r line; do
    if [[ $line == data:* ]]; then
        json_data="${line#data: }"
        echo "$json_data" | jq -r '.[] | "[\(.level)] \(.message)"' 2>/dev/null || true
    fi
done

# 6. Pause the spammer
echo -e "\n=== Pausing Spammer ==="
call_api "POST" "$BASE_URL/spammer/$spammer_id/pause" >/dev/null
echo "Spammer paused"

# 7. Get final logs
echo "=== Final Log Summary ==="
recent_logs=$(call_api "GET" "$BASE_URL/spammer/$spammer_id/logs")
log_count=$(echo "$recent_logs" | jq 'length')
echo "Total log entries: $log_count"

# Show last 5 log entries
echo "Last 5 log entries:"
echo "$recent_logs" | jq -r '.[-5:] | .[] | "[\(.level)] \(.time) - \(.message)"'
```

### Bulk Operations with Export/Import

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api"

# Export all current spammers
echo "=== Exporting Current Spammers ==="
curl -s -X POST "$BASE_URL/spammers/export" \
  -H "Content-Type: application/json" \
  -d '{}' > current_spammers.yaml

echo "Exported spammers to current_spammers.yaml"
cat current_spammers.yaml

# Modify configuration using yq (YAML processor)
echo -e "\n=== Modifying Configuration ==="
cp current_spammers.yaml modified_spammers.yaml

# Add "Modified" prefix to all spammer names
yq eval '.[] | .name = "Modified " + .name' -i modified_spammers.yaml

# Double throughput for all spammers that have it configured
yq eval '(.[] | select(.config.throughput) | .config.throughput) *= 2' -i modified_spammers.yaml

echo "Modified configuration:"
cat modified_spammers.yaml

# Import modified configuration
echo -e "\n=== Importing Modified Configuration ==="
modified_yaml=$(cat modified_spammers.yaml)
import_result=$(curl -s -X POST "$BASE_URL/spammers/import" \
  -H "Content-Type: application/json" \
  -d "{\"input\": $(echo "$modified_yaml" | jq -Rs .)}")

echo "Import result:"
echo "$import_result" | jq .

imported=$(echo "$import_result" | jq -r '.data.imported')
skipped=$(echo "$import_result" | jq -r '.data.skipped')
echo "Imported: $imported, Skipped: $skipped"

# Clean up temporary files
rm current_spammers.yaml modified_spammers.yaml
```

### Client Management

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api"

# Get all clients
echo "=== Current Client Status ==="
clients=$(curl -s "$BASE_URL/clients")
echo "$clients" | jq -r '.[] | "Client \(.index): \(.name) (\(.rpc_host)) - \(if .ready then "Ready" else "Not Ready" end)"'

echo -e "\n=== Updating Client Configuration ==="

# Process each client
echo "$clients" | jq -c '.[]' | while read -r client; do
    index=$(echo "$client" | jq -r '.index')
    name=$(echo "$client" | jq -r '.name')
    groups=$(echo "$client" | jq -r '.groups[]' 2>/dev/null)
    name_override=$(echo "$client" | jq -r '.name_override // empty')
    
    echo "Processing Client $index ($name)..."
    
    # Update client groups for load balancing
    if echo "$groups" | grep -q "mainnet"; then
        current_groups=$(echo "$client" | jq -r '.groups')
        new_groups=$(echo "$current_groups" | jq '. + ["load-balanced"] | unique')
        
        curl -s -X PUT "$BASE_URL/client/$index/group" \
            -H "Content-Type: application/json" \
            -d "{\"groups\": $new_groups}" > /dev/null
        
        echo "  → Added 'load-balanced' group"
    fi
    
    # Set custom names for easier identification
    if [[ -z "$name_override" ]]; then
        custom_name="Node-$((index + 1))-$name"
        curl -s -X PUT "$BASE_URL/client/$index/name" \
            -H "Content-Type: application/json" \
            -d "{\"name_override\": \"$custom_name\"}" > /dev/null
        
        echo "  → Set custom name: $custom_name"
    fi
done

echo -e "\n=== Updated Client Status ==="
updated_clients=$(curl -s "$BASE_URL/clients")
echo "$updated_clients" | jq -r '.[] | "Client \(.index): \(.name_override // .name) (\(.rpc_host)) - Groups: \(.groups | join(", "))"'
```

## SDKs and Libraries

While there are no official SDKs, the API is designed to be easily consumed by standard HTTP libraries:

### Recommended Tools

- **Bash/Shell**: `curl` for HTTP requests, `jq` for JSON processing, `yq` for YAML processing
- **JavaScript/Node.js**: `fetch`, `axios` for HTTP; `js-yaml` for configuration parsing
- **Go**: Standard `net/http` package or `resty`
- **Java**: `OkHttp`, `Apache HttpClient`
- **C#**: `HttpClient`

### OpenAPI/Swagger Integration

The API includes OpenAPI documentation that can be used to generate client SDKs:

1. Access the OpenAPI spec at `http://localhost:8080/docs/swagger.json`
2. Use tools like `swagger-codegen` or `openapi-generator` to generate clients
3. Example: `openapi-generator generate -i http://localhost:8080/docs/swagger.json -g python -o spamoor-python-client`

### Rate Limiting Considerations

While the API doesn't implement rate limiting, consider these best practices:

- **Polling**: Use reasonable intervals (1-5 seconds) when polling for status
- **Streaming**: Prefer SSE streaming for real-time updates over polling
- **Batch Operations**: Use export/import for bulk operations
- **Connection Reuse**: Use HTTP keep-alive and connection pooling

For high-frequency operations or production deployments, consider implementing client-side rate limiting and connection management.