#!/usr/bin/env bash
# =============================================================================
# 05-test.sh — End-to-end transfer tests (Cosmos<->EVM)
#
# Usage:
#   bash 05-test.sh                # Run both tests
#   bash 05-test.sh --cosmos-only  # Cosmos->EVM only
#   bash 05-test.sh --evm-only     # EVM->Cosmos only
#
# Required environment:
#   EVM_RELAYER_KEY (or HYP_KEY)  — for EVM->Cosmos transfer
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

EVM_RELAYER_KEY="${EVM_RELAYER_KEY:-${HYP_KEY:-}}"

# ─── Readiness Check ────────────────────────────────────────────────────────

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
  [[ "$ok" == "true" ]] || { log_err "Agents not ready. Run 04-agents.sh first."; exit 1; }

  # Wait for bitsong validator to be connected and syncing with the chain.
  # NOTE: Checkpoint index will be 0 until a message is dispatched (merkle tree is empty).
  # So we check validator logs for signs it's watching the chain, not for checkpoint production.
  log "Waiting for bitsong validator to connect to chain..."

  for i in $(seq 1 20); do
    local val_logs
    val_logs=$(docker logs hyperlane-validator-bitsong --tail 50 2>&1) || val_logs=""

    # Look for signs the validator is connected and watching the chain:
    # - "Watching for new messages" / "fetching" / "cursor" = actively scanning
    # - "Latest block" / "block_height" = synced with chain
    # - Any log output at all after startup = healthy
    if echo "$val_logs" | grep -qiE "watching|fetching|cursor|latest.block|block_height|merkle|tree_count|scanning"; then
      log_ok "Validator connected and watching chain"
      return 0
    fi

    # Also check if the validator has been running for >10s (it connects quickly)
    local started_at
    started_at=$(docker inspect hyperlane-validator-bitsong --format '{{.State.StartedAt}}' 2>/dev/null) || true
    if [[ -n "$started_at" ]]; then
      local start_epoch now_epoch
      start_epoch=$(date -d "$started_at" +%s 2>/dev/null) || start_epoch=0
      now_epoch=$(date +%s)
      local uptime=$(( now_epoch - start_epoch ))
      if [[ "$uptime" -gt 15 ]]; then
        log_ok "Validator running for ${uptime}s — proceeding"
        return 0
      fi
    fi

    if [[ $i -eq 1 ]]; then
      log "  Validator starting up — waiting for it to connect..."
    fi
    log "  [${i}/20] waiting..."
    sleep 3
  done

  log_warn "Could not confirm validator is syncing after 60s."
  log_warn "Continuing anyway — check: docker logs hyperlane-validator-bitsong"
}

# ─── Relay Diagnostics ──────────────────────────────────────────────────────

# Print a snapshot of what the relayer is doing (called during wait loops)
show_relay_progress() {
  local logs
  logs=$(docker logs hyperlane-relayer --tail 3000 2>&1 \
    | sed 's/\x1b\[[0-9;]*m//g') || return

  # Hyperlane v2 relayer log format:
  #   "Found log(s) in index range, ...cursor: RateLimitedContractSyncCursor { tip: N, ...
  #    domain: HyperlaneDomain(bitsong (7171)) }"
  #   "status: Finalized" for delivered txs
  #   "pool_size: N" in finality_stage

  local bitsong_tip basesep_tip finalized_count
  bitsong_tip=$(echo "$logs" | grep "HyperlaneDomain(bitsong" \
    | grep -oP 'tip: \K[0-9]+' | tail -1) || true
  basesep_tip=$(echo "$logs" | grep "HyperlaneDomain(basesepolia" \
    | grep -oP 'tip: \K[0-9]+' | tail -1) || true
  finalized_count=$(echo "$logs" | grep -ci "status: Finalized" 2>/dev/null) || finalized_count=0

  echo -n "    relayer: bitsong=${bitsong_tip:-?} basesep=${basesep_tip:-?} finalized=$finalized_count"

  # Check for errors (excluding CCIP noise)
  local real_errors
  real_errors=$(echo "$logs" | grep -i "error\|failed" \
    | grep -cv "0xa2827cb39\|CCIP\|verification" 2>/dev/null) || real_errors=0
  if [[ "$real_errors" -gt 0 ]]; then
    echo -n " ${RED}errors=$real_errors${NC}"
  fi
  echo
}

# ─── Cosmos → EVM ───────────────────────────────────────────────────────────

test_cosmos_to_evm() {
  log_step "Test: Cosmos -> EVM"

  local token_id evm_hyp_erc20 merkle_hook_id
  token_id=$(require_state "token_id" "token_id")
  evm_hyp_erc20=$(require_state "evm_hyp_erc20" "evm_hyp_erc20")
  merkle_hook_id=$(load_state "merkle_hook_id")

  # Send to EVM signer (so they have tokens for the EVM->Cosmos test)
  local evm_signer_addr evm_signer_bytes32
  evm_signer_addr=$(cast wallet address --private-key "$EVM_RELAYER_KEY" 2>/dev/null)
  [[ -n "$evm_signer_addr" ]] || { log_err "Cannot derive EVM signer address"; return 1; }
  local hex="${evm_signer_addr#0x}"; hex=$(echo "$hex" | tr '[:upper:]' '[:lower:]')
  evm_signer_bytes32=$(printf "0x%064s" "$hex" | tr ' ' '0')

  log "Sending 1000 ubtsg: BitSong -> Base Sepolia"
  log "  Token:     $token_id"
  log "  Recipient: $evm_signer_addr"

  local initial_supply
  initial_supply=$(cast call "$evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || initial_supply="0"
  log "  Initial totalSupply: $initial_supply"

  # Use merkle hook to bypass IGP payment (devnet)
  local tx_args=("$BINARY" tx warp transfer "$token_id" "$REMOTE_DOMAIN" "$evm_signer_bytes32" "1000"
    --max-hyperlane-fee "0ubtsg")
  [[ -n "$merkle_hook_id" ]] && tx_args+=(--custom-hook-id "$merkle_hook_id")

  submit_tx "Cosmos->EVM (1000 ubtsg)" "${tx_args[@]}"

  # Check if merkle tree count increased (confirms message was dispatched)
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
    current_supply=$(cast call "$evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || current_supply="0"
    if [[ "$current_supply" != "$initial_supply" ]]; then
      log_ok "Cosmos->EVM SUCCESS! totalSupply: $initial_supply -> $current_supply"
      save_state "cosmos_to_evm_test_passed" "true"
      return 0
    fi

    # Show diagnostics every 30s
    if (( i % 6 == 0 )); then
      echo -e "  ${CYAN}[${i}/60] totalSupply=$current_supply — checking relayer...${NC}"
      show_relay_progress

      # Also check validator checkpoint progress
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

  # Timed out — show diagnostics
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
  log_warn "  1. Check validator: docker logs hyperlane-validator-bitsong --tail 50"
  log_warn "  2. Check relayer:   docker logs hyperlane-relayer --tail 50"
  log_warn "  3. Run status:      bash status.sh"
  log_warn "  4. Retry test:      bash 05-test.sh --cosmos-only"
  return 1
}

# ─── EVM → Cosmos ───────────────────────────────────────────────────────────

test_evm_to_cosmos() {
  log_step "Test: EVM -> Cosmos"

  local evm_hyp_erc20 token_id
  evm_hyp_erc20=$(require_state "evm_hyp_erc20" "evm_hyp_erc20")
  token_id=$(require_state "token_id" "token_id")

  local cosmos_recipient_bytes32
  cosmos_recipient_bytes32=$(bech32_to_bytes32 "$VAL_ADDRESS")
  [[ -n "$cosmos_recipient_bytes32" ]] || { log_err "Failed to convert $VAL_ADDRESS to bytes32"; return 1; }

  # Check EVM balance first
  local evm_balance
  evm_balance=$(cast call "$evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || evm_balance="0"
  if [[ "$evm_balance" == "0" ]]; then
    log_err "HypERC20 totalSupply is 0 — no tokens to send back."
    log_err "Run Cosmos->EVM test first to mint tokens on EVM side."
    return 1
  fi

  log "Sending 1000 ubtsg: Base Sepolia -> BitSong"
  log "  HypERC20:  $evm_hyp_erc20 (balance: $evm_balance)"
  log "  Recipient: $VAL_ADDRESS"

  local initial_height
  initial_height=$("$BINARY" status --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
    jq -r '.sync_info.latest_block_height // "0"' 2>/dev/null) || initial_height="0"

  local evm_tx
  evm_tx=$(cast send "$evm_hyp_erc20" \
    "transferRemote(uint32,bytes32,uint256)" "$DOMAIN_ID" "$cosmos_recipient_bytes32" "1000" \
    --value 1 --private-key "$EVM_RELAYER_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true

  local evm_tx_hash
  evm_tx_hash=$(echo "$evm_tx" | jq -r '.transactionHash // empty' 2>/dev/null) || true
  [[ -n "$evm_tx_hash" ]] || { log_err "EVM tx failed"; echo "$evm_tx"; return 1; }

  log "EVM TX: $evm_tx_hash"
  log "Waiting for relayer delivery (timeout: 1500s — first run scans ~24M blocks)..."

  for i in $(seq 1 100); do
    sleep 15
    local events
    events=$("$BINARY" query txs \
      --query "coin_received.receiver='$VAL_ADDRESS'" \
      --limit 5 --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
      jq --arg h "$initial_height" '[.txs[]? | select((.height | tonumber) > ($h | tonumber))] | length' \
      2>/dev/null) || events="0"
    if [[ "${events:-0}" -gt 0 ]]; then
      log_ok "EVM->Cosmos SUCCESS! coin_received by $VAL_ADDRESS"
      save_state "evm_to_cosmos_test_passed" "true"
      return 0
    fi

    local height
    height=$("$BINARY" status --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
      jq -r '.sync_info.latest_block_height // "?"' 2>/dev/null) || height="?"

    # Show diagnostics every 60s
    if (( i % 4 == 0 )); then
      echo -e "  ${CYAN}[${i}/100] height=$height — checking relayer...${NC}"
      show_relay_progress
    else
      log "  [${i}/100] height=$height"
    fi
  done

  log_err "Timed out after 1500s!"
  log_warn "Last 10 relayer logs:"
  docker logs hyperlane-relayer --tail 20 2>&1 \
    | grep -v "0xa2827cb39\|CCIP Read" | tail -10 || true
  return 1
}

# ─── Summary ─────────────────────────────────────────────────────────────────

print_summary() {
  log_step "Full Summary"

  echo -e "${BOLD}Chain${NC}"
  echo "  ID:     $CHAIN_ID"
  echo "  RPC:    $NODE"
  echo "  Domain: $DOMAIN_ID"
  echo

  echo -e "${BOLD}Hyperlane${NC}"
  echo "  Mailbox:       $(load_state mailbox_id)"
  echo "  RoutingISM:    $(load_state routing_ism_id)"
  echo "  MerkleHook:    $(load_state merkle_hook_id)"
  echo "  IGP:           $(load_state igp_id)"
  echo "  Token:         $(load_state token_id)"
  echo

  echo -e "${BOLD}EVM (Base Sepolia)${NC}"
  echo "  HypERC20:      $(load_state evm_hyp_erc20)"
  echo "  MultisigISM:   $(load_state basesepolia_multisig_ism)"
  echo "  Validator:     $(load_state validator_addr)"
  echo

  echo -e "${BOLD}Transfer Tests${NC}"
  echo "  Cosmos->EVM:   $(load_state cosmos_to_evm_test_passed || echo 'not run')"
  echo "  EVM->Cosmos:   $(load_state evm_to_cosmos_test_passed || echo 'not run')"
  echo

  echo -e "${BOLD}Docker${NC}"
  echo "  docker logs hyperlane-validator-bitsong"
  echo "  docker logs hyperlane-validator-basesepolia"
  echo "  docker logs hyperlane-relayer"
  echo

  echo -e "${BOLD}State${NC}: $STATE_FILE"
}

# ─── Main ────────────────────────────────────────────────────────────────────

COSMOS_ONLY=false
EVM_ONLY=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --cosmos-only) COSMOS_ONLY=true; shift ;;
    --evm-only)    EVM_ONLY=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--cosmos-only|--evm-only]"
      echo "Required env: EVM_RELAYER_KEY (or HYP_KEY)"
      exit 0 ;;
    *) log_err "Unknown flag: $1"; exit 1 ;;
  esac
done

require_binary
require_jq
require_chain_running
command -v cast >/dev/null 2>&1 || { log_err "cast missing: foundryup"; exit 1; }
[[ -n "$EVM_RELAYER_KEY" ]] || { log_err "Set EVM_RELAYER_KEY or HYP_KEY"; exit 1; }

banner "Transfer Tests" "bitsong <-> basesepolia"

# Pre-flight: make sure agents are running + validator is caught up
check_agents_ready

if [[ "$EVM_ONLY" != "true" ]]; then
  test_cosmos_to_evm
fi

if [[ "$COSMOS_ONLY" != "true" ]]; then
  test_evm_to_cosmos
fi

print_summary
log_ok "Tests complete!"
