#!/usr/bin/env bash
# =============================================================================
# 07-fantoken-test.sh — End-to-end fantoken transfer tests (Cosmos<->EVM)
#
# Supports per-symbol testing. If --symbol is provided, loads state from
# per-symbol keys (ft_<sym>_*). Otherwise falls back to flat keys (ft_*).
#
# Usage:
#   bash 07-fantoken-test.sh --symbol clay                # Test clay route
#   bash 07-fantoken-test.sh --symbol clay --cosmos-only  # Cosmos->EVM only
#   bash 07-fantoken-test.sh                              # Test latest route
#
# Required environment:
#   EVM_RELAYER_KEY (or HYP_KEY)  — for EVM->Cosmos transfer
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

EVM_RELAYER_KEY="${EVM_RELAYER_KEY:-${HYP_KEY:-}}"

# =============================================================================
# Per-symbol state helpers (mirrors 06-fantoken-route.sh)
# =============================================================================

FT_KEY=""  # lowercased symbol, empty = use flat keys

ft_load_key() {
  if [[ -n "$FT_KEY" ]]; then
    load_state "ft_${FT_KEY}_$1"
  else
    load_state "ft_$1"
  fi
}

ft_save_key() {
  if [[ -n "$FT_KEY" ]]; then
    save_state "ft_${FT_KEY}_$1" "$2"
    # Also update flat aliases
    save_state "ft_$1" "$2"
  else
    save_state "ft_$1" "$2"
  fi
}

# =============================================================================
# Readiness Check
# =============================================================================

check_agents_ready() {
  log_step "Checking Agent Readiness"

  local ok=true
  for name in hyperlane-validator-bitsong hyperlane-validator-basesepolia hyperlane-relayer; do
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^${name}$"; then
      log_ok "$name running"
    else
      log_err "$name NOT running — start with: bash 04-agents.sh"
      ok=false
    fi
  done
  [[ "$ok" == "true" ]] || { log_err "Agents not ready. Run 04-agents.sh / 06-fantoken-route.sh first."; exit 1; }

  # Verify relayer has fantoken in whitelist
  local ft_restarted
  ft_restarted=$(ft_load_key "relayer_restarted")
  if [[ "$ft_restarted" != "true" ]]; then
    log_warn "Relayer may not have fantoken whitelist. Run 06-fantoken-route.sh first."
  fi
}

# =============================================================================
# Relay Diagnostics
# =============================================================================

show_relay_progress() {
  local logs
  logs=$(docker logs hyperlane-relayer --tail 3000 2>&1 \
    | sed 's/\x1b\[[0-9;]*m//g') || return

  local bitsong_tip basesep_tip finalized_count
  bitsong_tip=$(echo "$logs" | grep "HyperlaneDomain(bitsong" \
    | grep -oP 'tip: \K[0-9]+' | tail -1) || true
  basesep_tip=$(echo "$logs" | grep "HyperlaneDomain(basesepolia" \
    | grep -oP 'tip: \K[0-9]+' | tail -1) || true
  finalized_count=$(echo "$logs" | grep -ci "status: Finalized" 2>/dev/null) || finalized_count=0

  echo -n "    relayer: bitsong=${bitsong_tip:-?} basesep=${basesep_tip:-?} finalized=$finalized_count"

  local real_errors
  real_errors=$(echo "$logs" | grep -i "error\|failed" \
    | grep -cv "0xa2827cb39\|CCIP\|verification" 2>/dev/null) || real_errors=0
  if [[ "$real_errors" -gt 0 ]]; then
    echo -n " ${RED}errors=$real_errors${NC}"
  fi
  echo
}

# =============================================================================
# Cosmos -> EVM (fantoken)
# =============================================================================

test_ft_cosmos_to_evm() {
  local symbol_label="${FT_KEY:-latest}"
  log_step "Test: Cosmos -> EVM (Fantoken: $symbol_label)"

  local ft_token_id ft_evm_hyp_erc20 ft_denom merkle_hook_id
  ft_token_id=$(ft_load_key "token_id")
  ft_evm_hyp_erc20=$(ft_load_key "evm_hyp_erc20")
  ft_denom=$(ft_load_key "denom")
  merkle_hook_id=$(load_state "merkle_hook_id")

  [[ -n "$ft_token_id" ]] || { log_err "ft_token_id not in state"; return 1; }
  [[ -n "$ft_evm_hyp_erc20" ]] || { log_err "ft_evm_hyp_erc20 not in state"; return 1; }
  [[ -n "$ft_denom" ]] || { log_err "ft_denom not in state"; return 1; }

  # Send to EVM signer (so they have tokens for the EVM->Cosmos test)
  local evm_signer_addr evm_signer_bytes32
  evm_signer_addr=$(cast wallet address --private-key "$EVM_RELAYER_KEY" 2>/dev/null)
  [[ -n "$evm_signer_addr" ]] || { log_err "Cannot derive EVM signer address"; return 1; }
  local hex="${evm_signer_addr#0x}"; hex=$(echo "$hex" | tr '[:upper:]' '[:lower:]')
  evm_signer_bytes32=$(printf "0x%064s" "$hex" | tr ' ' '0')

  local transfer_amount=1000

  log "Sending $transfer_amount $ft_denom: BitSong -> Base Sepolia"
  log "  Token:     $ft_token_id"
  log "  Denom:     $ft_denom"
  log "  Recipient: $evm_signer_addr"

  local initial_supply
  initial_supply=$(cast call "$ft_evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || initial_supply="0"
  log "  Initial totalSupply: $initial_supply"

  # Use merkle hook to bypass IGP payment (devnet)
  local tx_args=("$BINARY" tx warp transfer "$ft_token_id" "$REMOTE_DOMAIN" "$evm_signer_bytes32" "$transfer_amount"
    --max-hyperlane-fee "0ubtsg")
  [[ -n "$merkle_hook_id" ]] && tx_args+=(--custom-hook-id "$merkle_hook_id")

  submit_tx "Cosmos->EVM ($transfer_amount $ft_denom)" "${tx_args[@]}"

  # Check if merkle tree count increased
  local merkle_count
  merkle_count=$("$BINARY" query hyperlane hooks merkle-tree-hooks \
    --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null \
    | jq -r '.merkle_tree_hooks[0].merkle_tree.count // "?"' 2>/dev/null) || merkle_count="?"
  log "  Merkle tree count after dispatch: $merkle_count"

  # Wait for relayer to deliver — 300s timeout with diagnostics
  log "Waiting for relayer to pick up, sign, and deliver (timeout: 300s)..."
  for i in $(seq 1 60); do
    sleep 5
    local current_supply
    current_supply=$(cast call "$ft_evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || current_supply="0"
    if [[ "$current_supply" != "$initial_supply" ]]; then
      log_ok "Cosmos->EVM (FT) SUCCESS! totalSupply: $initial_supply -> $current_supply"
      ft_save_key "cosmos_to_evm_test_passed" "true"
      return 0
    fi

    if (( i % 6 == 0 )); then
      echo -e "  ${CYAN}[${i}/60] totalSupply=$current_supply — checking relayer...${NC}"
      show_relay_progress

      local cp_file="$BITSONG_HOME/checkpoints-bitsong/index.json"
      if [[ -f "$cp_file" ]]; then
        local cp_idx
        cp_idx=$(cat "$cp_file" 2>/dev/null) || cp_idx="?"
        echo "    validator-bitsong: checkpoint_index=$cp_idx"
      fi
    else
      log "  [${i}/60] totalSupply=$current_supply"
    fi
  done

  log_err "Timed out after 300s!"
  log_warn "Diagnostic info:"
  log_warn "  Last 10 relayer log lines (filtered):"
  docker logs hyperlane-relayer --tail 20 2>&1 \
    | grep -v "0xa2827cb39\|CCIP Read" | tail -10 || true
  echo
  log_warn "  Validator checkpoint:"
  local cp_file="$BITSONG_HOME/checkpoints-bitsong/index.json"
  [[ -f "$cp_file" ]] && cat "$cp_file" || echo "  (no checkpoint file)"
  echo
  log_warn "Troubleshooting:"
  log_warn "  1. Check relayer whitelist includes ft_token_id"
  log_warn "  2. docker logs hyperlane-relayer --tail 50"
  log_warn "  3. bash status.sh"
  log_warn "  4. Retry: bash 07-fantoken-test.sh${FT_KEY:+ --symbol $FT_KEY} --cosmos-only"
  return 1
}

# =============================================================================
# EVM -> Cosmos (fantoken)
# =============================================================================

test_ft_evm_to_cosmos() {
  local symbol_label="${FT_KEY:-latest}"
  log_step "Test: EVM -> Cosmos (Fantoken: $symbol_label)"

  local ft_evm_hyp_erc20 ft_token_id ft_denom
  ft_evm_hyp_erc20=$(ft_load_key "evm_hyp_erc20")
  ft_token_id=$(ft_load_key "token_id")
  ft_denom=$(ft_load_key "denom")

  [[ -n "$ft_evm_hyp_erc20" ]] || { log_err "ft_evm_hyp_erc20 not in state"; return 1; }
  [[ -n "$ft_token_id" ]] || { log_err "ft_token_id not in state"; return 1; }
  [[ -n "$ft_denom" ]] || { log_err "ft_denom not in state"; return 1; }

  local cosmos_recipient_bytes32
  cosmos_recipient_bytes32=$(bech32_to_bytes32 "$VAL_ADDRESS")
  [[ -n "$cosmos_recipient_bytes32" ]] || { log_err "Failed to convert $VAL_ADDRESS to bytes32"; return 1; }

  # Check EVM balance first
  local evm_balance
  evm_balance=$(cast call "$ft_evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || evm_balance="0"
  if [[ "$evm_balance" == "0" ]]; then
    log_err "FT HypERC20 totalSupply is 0 — no tokens to send back."
    log_err "Run Cosmos->EVM fantoken test first."
    return 1
  fi

  local transfer_amount=1000

  # Get initial Cosmos balance for the fantoken
  local initial_balance
  initial_balance=$("$BINARY" query bank balance "$VAL_ADDRESS" "$ft_denom" \
    --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
    jq -r '.balance.amount // "0"' 2>/dev/null) || initial_balance="0"

  log "Sending $transfer_amount $ft_denom: Base Sepolia -> BitSong"
  log "  HypERC20:  $ft_evm_hyp_erc20 (supply: $evm_balance)"
  log "  Recipient: $VAL_ADDRESS"
  log "  Initial Cosmos balance: $initial_balance $ft_denom"

  local initial_height
  initial_height=$("$BINARY" status --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
    jq -r '.sync_info.latest_block_height // "0"' 2>/dev/null) || initial_height="0"

  local evm_tx
  evm_tx=$(cast send "$ft_evm_hyp_erc20" \
    "transferRemote(uint32,bytes32,uint256)" "$DOMAIN_ID" "$cosmos_recipient_bytes32" "$transfer_amount" \
    --value 1 --private-key "$EVM_RELAYER_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true

  local evm_tx_hash
  evm_tx_hash=$(echo "$evm_tx" | jq -r '.transactionHash // empty' 2>/dev/null) || true
  [[ -n "$evm_tx_hash" ]] || { log_err "EVM tx failed"; echo "$evm_tx"; return 1; }

  log "EVM TX: $evm_tx_hash"
  log "Waiting for relayer delivery (timeout: 1500s — first run scans ~24M blocks)..."

  for i in $(seq 1 100); do
    sleep 15
    local current_balance
    current_balance=$("$BINARY" query bank balance "$VAL_ADDRESS" "$ft_denom" \
      --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
      jq -r '.balance.amount // "0"' 2>/dev/null) || current_balance="0"

    if [[ "$current_balance" != "$initial_balance" ]]; then
      log_ok "EVM->Cosmos (FT) SUCCESS! Balance: $initial_balance -> $current_balance $ft_denom"
      ft_save_key "evm_to_cosmos_test_passed" "true"
      return 0
    fi

    local height
    height=$("$BINARY" status --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
      jq -r '.sync_info.latest_block_height // "?"' 2>/dev/null) || height="?"

    if (( i % 4 == 0 )); then
      echo -e "  ${CYAN}[${i}/100] height=$height balance=$current_balance — checking relayer...${NC}"
      show_relay_progress
    else
      log "  [${i}/100] height=$height balance=$current_balance"
    fi
  done

  log_err "Timed out after 1500s!"
  log_warn "Last 10 relayer logs:"
  docker logs hyperlane-relayer --tail 20 2>&1 \
    | grep -v "0xa2827cb39\|CCIP Read" | tail -10 || true
  return 1
}

# =============================================================================
# Summary
# =============================================================================

print_summary() {
  local symbol_label="${FT_KEY:-latest}"
  log_step "Fantoken Test Summary ($symbol_label)"

  echo -e "${BOLD}Fantoken${NC}"
  echo "  Denom:     $(ft_load_key denom)"
  echo "  Token ID:  $(ft_load_key token_id)"
  echo "  HypERC20:  $(ft_load_key evm_hyp_erc20)"
  echo

  echo -e "${BOLD}Transfer Tests${NC}"
  local c2e e2c
  c2e=$(ft_load_key "cosmos_to_evm_test_passed")
  e2c=$(ft_load_key "evm_to_cosmos_test_passed")
  echo "  Cosmos->EVM (FT):  ${c2e:-not run}"
  echo "  EVM->Cosmos (FT):  ${e2c:-not run}"
  echo
}

# =============================================================================
# Main
# =============================================================================

COSMOS_ONLY=false
EVM_ONLY=false
FT_SYMBOL=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --symbol)
      [[ -n "${2:-}" ]] || { log_err "--symbol requires a value"; exit 1; }
      FT_SYMBOL="$2"; shift 2 ;;
    --cosmos-only) COSMOS_ONLY=true; shift ;;
    --evm-only)    EVM_ONLY=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--symbol <sym>] [--cosmos-only|--evm-only]"
      echo ""
      echo "Options:"
      echo "  --symbol <sym>   Test a specific fantoken route (default: latest)"
      echo "  --cosmos-only    Cosmos->EVM only"
      echo "  --evm-only       EVM->Cosmos only"
      echo ""
      echo "Required env: EVM_RELAYER_KEY (or HYP_KEY)"
      exit 0 ;;
    *) log_err "Unknown flag: $1"; exit 1 ;;
  esac
done

if [[ -n "$FT_SYMBOL" ]]; then
  FT_KEY=$(echo "$FT_SYMBOL" | tr '[:upper:]' '[:lower:]')
fi

require_binary
require_jq
require_chain_running
command -v cast >/dev/null 2>&1 || { log_err "cast missing: foundryup"; exit 1; }
[[ -n "$EVM_RELAYER_KEY" ]] || { log_err "Set EVM_RELAYER_KEY or HYP_KEY"; exit 1; }

if [[ -n "$FT_KEY" ]]; then
  banner "Fantoken Transfer Tests" "$FT_SYMBOL — bitsong <-> basesepolia"
else
  banner "Fantoken Transfer Tests" "latest route — bitsong <-> basesepolia"
fi

check_agents_ready

if [[ "$EVM_ONLY" != "true" ]]; then
  test_ft_cosmos_to_evm
fi

if [[ "$COSMOS_ONLY" != "true" ]]; then
  test_ft_evm_to_cosmos
fi

print_summary
log_ok "Fantoken tests complete!"
