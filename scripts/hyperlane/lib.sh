#!/usr/bin/env bash
# =============================================================================
# lib.sh — Shared library for BitSong Hyperlane scripts
#
# Source this file from any script:
#   source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"
# =============================================================================

[[ -n "${_HYPERLANE_LIB_LOADED:-}" ]] && return 0
_HYPERLANE_LIB_LOADED=1

set -euo pipefail

# =============================================================================
# Paths
# =============================================================================

_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Walk up to project root (contains go.mod)
_PROJECT_DIR="$_LIB_DIR"
while [[ "$_PROJECT_DIR" != "/" && ! -f "$_PROJECT_DIR/go.mod" ]]; do
  _PROJECT_DIR="$(dirname "$_PROJECT_DIR")"
done
if [[ ! -f "$_PROJECT_DIR/go.mod" ]]; then
  echo "ERROR: Could not find project root (go.mod)." >&2
  exit 1
fi

# Auto-load .env if present (private keys, etc. — gitignored)
if [[ -f "$_LIB_DIR/.env" ]]; then
  set -a; source "$_LIB_DIR/.env"; set +a
fi

# =============================================================================
# Colors & Logging
# =============================================================================

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

log()      { echo -e "${CYAN}[$(date '+%H:%M:%S')]${NC} $*"; }
log_ok()   { echo -e "${GREEN}[$(date '+%H:%M:%S')] ✓${NC} $*"; }
log_warn() { echo -e "${YELLOW}[$(date '+%H:%M:%S')] !${NC} $*"; }
log_err()  { echo -e "${RED}[$(date '+%H:%M:%S')] ✗${NC} $*"; }
log_step() { echo -e "\n${BOLD}${CYAN}═══ $* ═══${NC}\n"; }

banner() {
  echo -e "${BOLD}${CYAN}"
  echo "  ╔══════════════════════════════════════════════╗"
  printf "  ║   %-43s║\n" "$1"
  printf "  ║   %-43s║\n" "$2"
  echo "  ╚══════════════════════════════════════════════╝"
  echo -e "${NC}"
}

# =============================================================================
# Chain Configuration (all overridable via environment)
# =============================================================================

CHAIN_ID="${CHAIN_ID:-localbitsong-hyperlane}"
BITSONG_HOME="${BITSONG_HOME:-$HOME/.localbitsong-hyperlane}"
MONIKER="${MONIKER:-val}"
KEY_NAME="${KEY_NAME:-val}"
KEYRING_BACKEND="${KEYRING_BACKEND:-test}"
DENOM="${DENOM:-ubtsg}"
NODE="${NODE:-tcp://localhost:26657}"

# Auto-detect bitsongd
if [[ -z "${BINARY:-}" ]]; then
  if [[ -x "$_PROJECT_DIR/build/bitsongd" ]]; then
    BINARY="$_PROJECT_DIR/build/bitsongd"
  elif [[ -x "./build/bitsongd" ]]; then
    BINARY="./build/bitsongd"
  else
    BINARY="bitsongd"
  fi
fi

# Validator mnemonic — loaded from .env (never hardcoded)
# Required by 01-chain.sh; optional for other scripts.
VAL_MNEMONIC="${VAL_MNEMONIC:-}"

# Derive VAL_ADDRESS from keyring (coin type 639, set by bitsongd)
VAL_ADDRESS="${VAL_ADDRESS:-}"
if [[ -z "$VAL_ADDRESS" && -d "$BITSONG_HOME/keyring-$KEYRING_BACKEND" ]]; then
  VAL_ADDRESS=$("$BINARY" keys show "$KEY_NAME" \
    --keyring-backend "$KEYRING_BACKEND" --home "$BITSONG_HOME" -a 2>/dev/null) || true
fi

# =============================================================================
# Hyperlane / EVM Configuration
# =============================================================================

DOMAIN_ID="${DOMAIN_ID:-7171}"
REMOTE_DOMAIN="${REMOTE_DOMAIN:-84532}"
BASESEPOLIA_MAILBOX="${BASESEPOLIA_MAILBOX:-0x6966b0E55883d49BFB24539356a2f8A673E02039}"
EVM_RPC="${EVM_RPC:-https://sepolia.base.org}"

VALIDATOR_ADDRESSES="${VALIDATOR_ADDRESSES:-0x1111111111111111111111111111111111111111,0x2222222222222222222222222222222222222222,0x3333333333333333333333333333333333333333}"
VALIDATOR_THRESHOLD="${VALIDATOR_THRESHOLD:-2}"
TOKEN_EXCHANGE_RATE="${TOKEN_EXCHANGE_RATE:-10000000000}"
GAS_PRICE="${GAS_PRICE:-1000000000}"
GAS_OVERHEAD="${GAS_OVERHEAD:-75000}"
ENROLL_GAS="${ENROLL_GAS:-300000}"

# Docker agents
DOCKER_IMAGE="${DOCKER_IMAGE:-gcr.io/abacus-labs-dev/hyperlane-agent:agents-v2.0.0}"
BASESEPOLIA_ISM_FACTORY="${BASESEPOLIA_ISM_FACTORY:-0xfc6e546510dC9d76057F1f76633FCFfC188CB213}"
VALIDATOR_ANNOUNCE_CONTRACT="${VALIDATOR_ANNOUNCE_CONTRACT:-0x20c44b1E3BeaDA1e9826CFd48BeEDABeE9871cE9}"

# State file (shared across all phases)
STATE_FILE="$BITSONG_HOME/hyperlane-state.json"

# =============================================================================
# State Management
# =============================================================================

save_state() {
  local key="$1" value="$2"
  [[ ! -f "$STATE_FILE" ]] && echo '{}' > "$STATE_FILE"
  local tmp="${STATE_FILE}.tmp"
  jq --arg k "$key" --arg v "$value" '.[$k] = $v' "$STATE_FILE" > "$tmp" && mv "$tmp" "$STATE_FILE"
}

load_state() {
  local key="$1"
  [[ -f "$STATE_FILE" ]] && jq -r --arg k "$key" '.[$k] // empty' "$STATE_FILE" 2>/dev/null || true
}

# =============================================================================
# Chain Helpers
# =============================================================================

wait_for_block() {
  log "Waiting for chain to produce blocks..."
  for _ in $(seq 1 60); do
    local height
    height=$("$BINARY" status --node "$NODE" --home "$BITSONG_HOME" 2>&1 \
      | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // "0"' 2>/dev/null) || height="0"
    if [[ "$height" =~ ^[0-9]+$ && "$height" -gt 0 ]]; then
      log_ok "Chain producing blocks (height=$height)"
      return 0
    fi
    sleep 2
  done
  log_err "Chain not producing blocks after 120s — check: $BITSONG_HOME/chain.log"
  return 1
}

# =============================================================================
# Transaction Helpers
# =============================================================================

TX_RESULT=""

# Submit a bitsongd TX and wait for inclusion (up to 60s).
# Usage: submit_tx "description" $BINARY tx hyperlane ...
submit_tx() {
  local desc="$1"; shift
  TX_RESULT=""
  log "Submitting: $desc"

  local result
  result=$("$@" \
    --from "$KEY_NAME" --keyring-backend "$KEYRING_BACKEND" \
    --chain-id "$CHAIN_ID" --node "$NODE" \
    --gas auto --gas-adjustment 1.5 --fees "10000${DENOM}" \
    --output json -y --home "$BITSONG_HOME" 2>&1) || true

  # --gas auto prints "gas estimate: N" before JSON
  local json_result
  json_result=$(echo "$result" | grep '^\{' | tail -1) || true

  local txhash
  txhash=$(echo "$json_result" | jq -r '.txhash // empty' 2>/dev/null) || true
  if [[ -z "$txhash" ]]; then
    log_err "Failed to submit tx. Output:"; echo "$result"; return 1
  fi

  log "  TX: $txhash — waiting for inclusion..."
  for _ in $(seq 1 30); do
    sleep 2
    TX_RESULT=$("$BINARY" query tx "$txhash" --output json --node "$NODE" --home "$BITSONG_HOME" 2>&1) || true
    if echo "$TX_RESULT" | jq -e '.code != null' >/dev/null 2>&1; then break; fi
    TX_RESULT=""
  done

  [[ -n "$TX_RESULT" ]] || { log_err "TX not found after 60s"; return 1; }

  local code
  code=$(echo "$TX_RESULT" | jq -r '.code')
  if [[ "$code" != "0" ]]; then
    log_err "TX failed (code=$code)"
    echo "$TX_RESULT" | jq -r '.raw_log // "unknown"' 2>/dev/null || true
    return 1
  fi
  log_ok "TX succeeded (code=0)"
}

# Extract hex ID (0x...) from the last TX_RESULT.
extract_id() {
  local id=""

  # Strategy 1: msg_responses
  id=$(echo "$TX_RESULT" | jq -r '.msg_responses[0].id // empty' 2>/dev/null) || true
  [[ -n "$id" && "$id" != "null" ]] && { echo "$id"; return; }

  # Strategy 2: decode hex data field (protobuf response)
  local data
  data=$(echo "$TX_RESULT" | jq -r '.data // empty' 2>/dev/null) || true
  if [[ -n "$data" ]]; then
    id=$(echo "$data" | xxd -r -p | grep -aoP '0x[0-9a-f]{64}' | tail -1) || true
    [[ -n "$id" ]] && { echo "$id"; return; }
  fi

  # Strategy 3: Hyperlane/Warp events
  id=$(echo "$TX_RESULT" | jq -r '
    [.events[]? | select(.type | test("hyperlane|warp")) |
     select(.type | test("[Cc]reate")) |
     .attributes[]? | select(.key | test("_id$|^id$")) |
     select(.key != "msg_index") | .value
    ] | first // empty' 2>/dev/null) || true
  [[ -n "$id" && "$id" != "null" ]] && { echo "$id" | sed 's/^"//;s/"$//'; return; }

  echo ""
}

# =============================================================================
# Address Conversion
# =============================================================================

# 20-byte EVM address → 32-byte Hyperlane hex (left-padded)
evm_to_bytes32() {
  local addr="${1#0x}"
  addr=$(echo "$addr" | tr '[:upper:]' '[:lower:]')
  printf "0x%064s" "$addr" | tr ' ' '0'
}

# bech32 Cosmos address → 32-byte hex
bech32_to_bytes32() {
  local bech32_addr="$1" hex=""

  hex=$("$BINARY" keys parse "$bech32_addr" --output json 2>/dev/null \
    | jq -r '.bytes // empty' 2>/dev/null) || true

  if [[ -z "$hex" ]] && command -v python3 >/dev/null 2>&1; then
    hex=$(python3 -c "
CHARSET='qpzry9x8gf2tvdw0s3jn54khce6mua7l'
def decode(addr):
    _, dp = addr.rsplit('1', 1)
    v = [CHARSET.index(c) for c in dp][:-6]
    acc, bits, out = 0, 0, []
    for d in v:
        acc = (acc << 5) | d; bits += 5
        while bits >= 8: bits -= 8; out.append((acc >> bits) & 0xff)
    return bytes(out).hex()
print(decode('$bech32_addr'))
" 2>/dev/null) || true
  fi

  [[ -z "$hex" ]] && { log_err "Failed to decode bech32: $bech32_addr"; return 1; }
  printf "0x%064s" "$hex" | tr ' ' '0'
}

# =============================================================================
# Preflight Helpers
# =============================================================================

require_binary()        { [[ -x "$BINARY" ]] || command -v "$BINARY" >/dev/null 2>&1 || { log_err "bitsongd not found. Build: LEDGER_ENABLED=false make build"; exit 1; }; }
require_jq()            { command -v jq >/dev/null 2>&1 || { log_err "jq missing: sudo apt install jq"; exit 1; }; }
require_chain_running() { "$BINARY" status --node "$NODE" --home "$BITSONG_HOME" >/dev/null 2>&1 || { log_err "Chain not running at $NODE"; exit 1; }; }

require_state() {
  local key="$1" label="${2:-$1}" val
  val=$(load_state "$key")
  [[ -n "$val" ]] || { log_err "$label not in state file. Run previous phase first."; exit 1; }
  echo "$val"
}
