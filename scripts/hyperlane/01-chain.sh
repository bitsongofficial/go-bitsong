#!/usr/bin/env bash
# =============================================================================
# 01-chain.sh — Initialize and start a single-validator BitSong localnet
#
# Usage:
#   bash 01-chain.sh           # Init (if needed) + start
#   bash 01-chain.sh --clean   # Wipe state and re-init
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

# Test accounts (from tests/localbitsong/scripts/setup.sh)
TEST_ADDRESSES=(
  "bitsong1regz7kj3ylg2dn9rl8vwrhclkgz528mf0tfsck"
  "bitsong1hvrhhex6wfxh7r7nnc3y39p0qlmff6v9t5rc25"
  "bitsong175vgzztymvvcxvqun54nlu9dq6856thgvyl5sa"
  "bitsong1t8nznzj4sd6zzutwdmslgy4dcxyd2jafz7822x"
  "bitsong14vdrvstsffj8mq5e4fhm6y2hpfxtedajczsj5d"
  "bitsong1vwe5hay74v0vhuzdhadteyqfasu5d7tdf83pyy"
  "bitsong16866dezn6ez2qpmpcrrv9cyud8v8c7ufnzwhhh"
  "bitsong1tlwh75lvu35nw9vcg557mxhspz5s88t6vzscd8"
  "bitsong16z9wj8n5f3zgzwspw0r9sj9v7k7hdasqj95us9"
  "bitsong1gulaxnca7rped0grw0lz4h4zy0xn3ttvmlad8x"
)
GENESIS_AMOUNT="10000000000000${DENOM}"
STAKE_AMOUNT="5000000000000${DENOM}"

# ─── Init ────────────────────────────────────────────────────────────────────

init_chain() {
  log_step "Initializing Chain"

  log "Initializing chain (chain-id=$CHAIN_ID)..."
  echo "$VAL_MNEMONIC" | "$BINARY" init "$MONIKER" \
    --chain-id="$CHAIN_ID" --home="$BITSONG_HOME" --recover -o > /dev/null 2>&1

  # Genesis modifications
  log "Modifying genesis.json..."
  local genesis="$BITSONG_HOME/config/genesis.json"
  local tmp="${genesis}.tmp"
  jq '
    .app_state.staking.params.bond_denom = "ubtsg" |
    .app_state.staking.params.unbonding_time = "240s" |
    .app_state.crisis.constant_fee.denom = "ubtsg" |
    .app_state.gov.params.voting_period = "60s" |
    .app_state.gov.params.expedited_voting_period = "30s" |
    .app_state.gov.params.min_deposit[0].denom = "ubtsg" |
    .app_state.gov.params.expedited_min_deposit[0].denom = "ubtsg" |
    .app_state.mint.params.mint_denom = "ubtsg" |
    .app_state.bank.denom_metadata = [{
      "description": "Registered denom ubtsg for localbitsong-hyperlane testing",
      "denom_units": [{"denom": "ubtsg", "exponent": 0}],
      "base": "ubtsg", "display": "ubtsg", "name": "ubtsg", "symbol": "ubtsg"
    }]
  ' "$genesis" > "$tmp" && mv "$tmp" "$genesis"
  log_ok "Genesis modified"

  # Keys + accounts
  log "Adding validator key (coin-type 639)..."
  echo "$VAL_MNEMONIC" | "$BINARY" keys add "$KEY_NAME" \
    --recover --coin-type 639 --keyring-backend="$KEYRING_BACKEND" --home="$BITSONG_HOME" > /dev/null 2>&1

  # Derive VAL_ADDRESS from the just-created key
  VAL_ADDRESS=$("$BINARY" keys show "$KEY_NAME" \
    --keyring-backend "$KEYRING_BACKEND" --home "$BITSONG_HOME" -a 2>/dev/null)
  [[ -n "$VAL_ADDRESS" ]] || { log_err "Failed to derive VAL_ADDRESS from keyring"; return 1; }
  log_ok "Validator address: $VAL_ADDRESS"

  log "Adding genesis accounts..."
  "$BINARY" genesis add-genesis-account "$VAL_ADDRESS" "$GENESIS_AMOUNT" \
    --home="$BITSONG_HOME" > /dev/null 2>&1
  for addr in "${TEST_ADDRESSES[@]}"; do
    "$BINARY" genesis add-genesis-account "$addr" "$GENESIS_AMOUNT" \
      --home="$BITSONG_HOME" > /dev/null 2>&1
  done
  log_ok "Added $(( ${#TEST_ADDRESSES[@]} + 1 )) genesis accounts"

  # Gentx
  log "Creating gentx..."
  "$BINARY" genesis gentx "$KEY_NAME" "$STAKE_AMOUNT" \
    --keyring-backend="$KEYRING_BACKEND" --chain-id="$CHAIN_ID" \
    --home="$BITSONG_HOME" > /dev/null 2>&1
  "$BINARY" genesis collect-gentxs --home="$BITSONG_HOME" > /dev/null 2>&1
  "$BINARY" genesis validate-genesis --home="$BITSONG_HOME" > /dev/null 2>&1
  log_ok "Genesis validated"

  # Config tweaks
  log "Modifying config.toml + app.toml..."
  local config="$BITSONG_HOME/config/config.toml"
  local app_toml="$BITSONG_HOME/config/app.toml"
  sed -i 's/seeds = ".*"/seeds = ""/' "$config"
  sed -i 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|' "$config"
  sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/' "$config"
  sed -i '/^\[api\]/,/^\[/{s/^enable = false/enable = true/}' "$app_toml"
  sed -i 's/^swagger = false/swagger = true/' "$app_toml"
  sed -i 's/^enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "$app_toml"
  sed -i 's/^minimum-gas-prices = ".*"/minimum-gas-prices = "0ubtsg"/' "$app_toml"
  log_ok "Config files modified"
}

# ─── Start ───────────────────────────────────────────────────────────────────

start_chain() {
  log_step "Starting Chain"

  if [[ -f "$BITSONG_HOME/chain.pid" ]]; then
    local pid; pid=$(cat "$BITSONG_HOME/chain.pid")
    if kill -0 "$pid" 2>/dev/null; then
      log_ok "Chain already running (PID=$pid)"
      wait_for_block; return 0
    fi
    log_warn "Stale PID file, removing"
    rm -f "$BITSONG_HOME/chain.pid"
  fi

  log "Starting bitsongd (logs: $BITSONG_HOME/chain.log)..."
  "$BINARY" start --home "$BITSONG_HOME" > "$BITSONG_HOME/chain.log" 2>&1 &
  local chain_pid=$!; disown "$chain_pid"
  echo "$chain_pid" > "$BITSONG_HOME/chain.pid"

  sleep 3
  if ! kill -0 "$chain_pid" 2>/dev/null; then
    log_err "Chain died immediately. Last 20 lines:"
    tail -20 "$BITSONG_HOME/chain.log" || true
    return 1
  fi

  log_ok "Chain started (PID=$chain_pid)"
  wait_for_block
}

# ─── Main ────────────────────────────────────────────────────────────────────

CLEAN=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --clean) CLEAN=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--clean]"
      echo "  --clean   Remove existing chain data before starting"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

require_binary
require_jq
[[ -n "$VAL_MNEMONIC" ]] || { log_err "Set VAL_MNEMONIC in scripts/hyperlane/.env"; exit 1; }

banner "BitSong Localnet (Hyperlane)" "Chain: $CHAIN_ID"

if [[ "$CLEAN" == "true" ]]; then
  log "Cleaning previous state at $BITSONG_HOME..."
  if [[ -f "$BITSONG_HOME/chain.pid" ]]; then
    kill "$(cat "$BITSONG_HOME/chain.pid")" 2>/dev/null || true; sleep 2
  fi
  rm -rf "$BITSONG_HOME"
  log_ok "Clean complete"
fi

if [[ ! -d "$BITSONG_HOME/config" ]]; then
  init_chain
else
  log_ok "Chain already initialized at $BITSONG_HOME"
fi

start_chain
log_ok "Chain ready!"
