#!/usr/bin/env bash
# =============================================================================
# status.sh — Dashboard for Hyperlane agent sync status
#
# Shows: container health, validator checkpoints, relayer block indexing,
#        message processing, and warp token state.
#
# Usage:
#   bash status.sh              # One-shot status
#   bash status.sh --watch      # Refresh every 10s
#   bash status.sh --watch 5    # Refresh every 5s
#   bash status.sh --logs       # Show recent relayer logs
#   bash status.sh --logs 50    # Show last 50 lines
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

# ─── Docker Status ───────────────────────────────────────────────────────────

print_containers() {
  echo -e "${BOLD}Docker Containers${NC}"
  echo -e "  NAME                                STATUS     UPTIME          HEALTH"
  echo -e "  ─────────────────────────────────── ────────── ─────────────── ──────"

  for name in hyperlane-validator-bitsong hyperlane-validator-basesepolia hyperlane-relayer; do
    local info
    info=$(docker inspect "$name" --format '{{.State.Status}}|{{.State.StartedAt}}|{{if .State.Health}}{{.State.Health.Status}}{{else}}n/a{{end}}' 2>/dev/null) || info=""
    if [[ -z "$info" ]]; then
      echo -e "  ${name}$(printf '%*s' $((35 - ${#name})) '') ${RED}absent${NC}     -               -"
    else
      local status started health
      IFS='|' read -r status started health <<< "$info"
      local uptime="-"
      if [[ "$status" == "running" && -n "$started" ]]; then
        local start_epoch now_epoch
        start_epoch=$(date -d "$started" +%s 2>/dev/null) || start_epoch=0
        now_epoch=$(date +%s)
        local diff=$(( now_epoch - start_epoch ))
        if [[ $diff -lt 60 ]]; then uptime="${diff}s"
        elif [[ $diff -lt 3600 ]]; then uptime="$(( diff / 60 ))m$(( diff % 60 ))s"
        else uptime="$(( diff / 3600 ))h$(( (diff % 3600) / 60 ))m"
        fi
      fi
      local color="$GREEN"
      [[ "$status" != "running" ]] && color="$RED"
      echo -e "  ${name}$(printf '%*s' $((35 - ${#name})) '') ${color}${status}${NC}$(printf '%*s' $((10 - ${#status})) '') ${uptime}$(printf '%*s' $((15 - ${#uptime})) '') ${health:-n/a}"
    fi
  done
  echo
}

# ─── Chain Heights ───────────────────────────────────────────────────────────

print_chain_heights() {
  echo -e "${BOLD}Chain Heights${NC}"

  local bitsong_height="-"
  if "$BINARY" status --node "$NODE" --home "$BITSONG_HOME" >/dev/null 2>&1; then
    bitsong_height=$("$BINARY" status --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null \
      | jq -r '.sync_info.latest_block_height // "?"' 2>/dev/null) || bitsong_height="?"
  fi
  echo -e "  BitSong (local):    $bitsong_height"

  if command -v cast >/dev/null 2>&1; then
    local basesep_height
    basesep_height=$(cast block-number --rpc-url "$EVM_RPC" 2>/dev/null) || basesep_height="?"
    echo -e "  Base Sepolia:       $basesep_height"
  fi
  echo
}

# ─── Validator Checkpoints ───────────────────────────────────────────────────

print_checkpoints() {
  echo -e "${BOLD}Validator Checkpoints${NC}"

  for chain in bitsong basesepolia; do
    local cp_dir="$BITSONG_HOME/checkpoints-${chain}"
    if [[ ! -d "$cp_dir" ]]; then
      echo -e "  $chain: ${YELLOW}no checkpoint dir${NC}"
      continue
    fi

    local latest_idx="-"
    local checkpoint_file="$cp_dir/index.json"
    if [[ -f "$checkpoint_file" ]]; then
      # index.json contains either a plain number or JSON with .latest_index/.index
      local raw
      raw=$(cat "$checkpoint_file" 2>/dev/null) || raw=""
      if [[ "$raw" =~ ^[0-9]+$ ]]; then
        latest_idx="$raw"
      else
        latest_idx=$(echo "$raw" | jq -r '.latest_index // .index // "?"' 2>/dev/null) || latest_idx="?"
      fi
    fi

    local sig_count
    sig_count=$(find "$cp_dir" -name "*.json" -not -name "checkpoint_*" -not -name "announcement*" 2>/dev/null | wc -l) || sig_count=0

    echo -e "  $chain: latest_index=$latest_idx  signatures=$sig_count"
  done
  echo
}

# ─── Relayer Sync Progress ──────────────────────────────────────────────────

print_relayer_sync() {
  echo -e "${BOLD}Relayer Sync Progress${NC}"

  if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^hyperlane-relayer$"; then
    echo -e "  ${YELLOW}Relayer not running${NC}"
    echo
    return
  fi

  # Strip ANSI helper — Rust tracing embeds color codes around field names
  local strip_ansi='s/\x1b\[[0-9;]*m//g'

  # Per-chain block indexing — single docker call, then split locally.
  # Use --since 5m (not --tail) because basesepolia tx_id lookups are extremely
  # verbose (thousands of lines between each ~30s cursor entry).
  #
  # Two cursor types exist:
  #   RateLimitedContractSyncCursor — has tip/next_block/start_block (gas_payments)
  #   ForwardBackwardSequenceAwareSyncCursor — has at_block/sequence (dispatched_messages)
  # We prefer RateLimited (more useful block progress data), fall back to any.
  local all_sync_logs
  all_sync_logs=$(docker logs hyperlane-relayer --since 5m 2>&1 \
    | sed "$strip_ansi" \
    | grep "HyperlaneDomain(") || all_sync_logs=""

  for chain in bitsong basesepolia; do
    # Prefer RateLimited cursor (has tip/next_block — best for block progress)
    local chain_log
    chain_log=$(echo "$all_sync_logs" \
      | grep "HyperlaneDomain(${chain}" \
      | grep "RateLimitedContractSyncCursor" \
      | tail -1) || true
    # Fallback: any cursor type
    if [[ -z "$chain_log" ]]; then
      chain_log=$(echo "$all_sync_logs" \
        | grep "HyperlaneDomain(${chain}" \
        | tail -1) || true
    fi

    if [[ -n "$chain_log" ]]; then
      local sync_status
      sync_status=$(echo "$chain_log" | grep -oP 'estimated_time_to_sync: "\K[^"]+') || sync_status=""

      # Try RateLimited fields: tip, next_block, start_block, chunk_size
      local tip next_block start_block chunk_size
      tip=$(echo "$chain_log" | grep -oP 'tip: \K[0-9]+' | head -1) || tip=""
      next_block=$(echo "$chain_log" | grep -oP 'next_block: \K[0-9]+' | head -1) || next_block=""
      start_block=$(echo "$chain_log" | grep -oP 'start_block: \K[0-9]+' | head -1) || start_block=""
      chunk_size=$(echo "$chain_log" | grep -oP 'chunk_size: \K[0-9]+' | head -1) || chunk_size=""

      # Fallback for ForwardBackward cursor: use at_block from last_indexed_snapshot
      if [[ -z "$tip" ]]; then
        tip=$(echo "$chain_log" | grep -oP 'at_block: \K[0-9]+' | tail -1) || tip=""
        # Use range field: "range: A..=B"
        next_block=$(echo "$chain_log" | grep -oP 'range: \K[0-9]+') || next_block=""
      fi

      if [[ -n "$tip" ]]; then
        # Calculate progress percentage
        local pct="-"
        if [[ -n "$start_block" && -n "$next_block" ]]; then
          local scanned=$(( next_block - start_block ))
          local total=$(( tip - start_block + 1 ))
          if [[ "$total" -gt 0 ]]; then
            pct=$(( scanned * 100 / total ))
          fi
        fi

        if [[ "$sync_status" == "synced" ]]; then
          echo -e "  $chain: block ${tip}  ${GREEN}synced${NC}"
        else
          # Not synced — show block progress, percentage, chunk size, ETA
          local status_text=""
          if [[ -n "$sync_status" ]]; then
            status_text="ETA: $sync_status"
          else
            status_text="scanning"
          fi

          local progress_detail="${next_block:-?}/${tip}"
          if [[ "$pct" != "-" ]]; then
            progress_detail="${progress_detail} (${pct}%)"
          fi
          echo -e "  $chain: ${progress_detail}  chunk=${chunk_size:-?}  ${YELLOW}${status_text}${NC}"
        fi
      else
        echo -e "  $chain: ${YELLOW}syncing...${NC}"
      fi
    else
      echo -e "  $chain: ${YELLOW}no progress data yet${NC}"
    fi
  done
  echo

  # Grab recent logs for stats (submitted/finalized/errors)
  local stats_logs
  stats_logs=$(docker logs hyperlane-relayer --since 10m 2>&1 \
    | sed "$strip_ansi") || stats_logs=""

  # Pending tx_id lookups (basesepolia generates many — shows indexer backlog)
  local pending_ids
  pending_ids=$(echo "$stats_logs" | grep -oP 'pending_ids: \K[0-9]+' | tail -1) || true
  if [[ -n "$pending_ids" && "$pending_ids" != "0" ]]; then
    echo -e "  Pending tx_id lookups: $pending_ids"
  fi

  # Finality pool size
  local pool_size
  pool_size=$(echo "$stats_logs" | grep -oP 'pool_size: \K[0-9]+' | tail -1) || true
  if [[ -n "$pool_size" ]]; then
    echo -e "  Finality pool: $pool_size txs"
  fi

  # Relay stats — match actual v2 log patterns
  # Transaction lifecycle: submitting → PendingInclusion → Mempool → Included → Finalized
  # "Message successfully processed" = confirmed delivery
  local submitted_count finalized_count processed_count error_count real_errors
  submitted_count=$(echo "$stats_logs" | grep -ci "submitting transaction" 2>/dev/null) || submitted_count=0
  finalized_count=$(echo "$stats_logs" | grep -ci "new_status: Finalized" 2>/dev/null) || finalized_count=0
  processed_count=$(echo "$stats_logs" | grep -ci "Message successfully processed" 2>/dev/null) || processed_count=0
  error_count=$(echo "$stats_logs" | grep -ci "error\|failed" 2>/dev/null) || error_count=0
  real_errors=$(echo "$stats_logs" | grep -i "error\|failed" | grep -cv "0xa2827cb39\|CCIP\|verification" 2>/dev/null) || real_errors=0

  echo -e "  Messages:  ${GREEN}${processed_count} delivered${NC}  ${submitted_count} submitted  ${finalized_count} finalized"
  if [[ "$real_errors" -gt 0 ]]; then
    echo -e "  Errors:    ${RED}${real_errors} real${NC} (${error_count} total incl. CCIP noise)"
  else
    echo -e "  Errors:    0 (${error_count} total incl. CCIP noise)"
  fi

  # Show latest relay transaction status (most useful when waiting for a transfer)
  local last_tx_status
  last_tx_status=$(echo "$stats_logs" \
    | grep -oP 'Updating tx status.*?new_status: \K\w+' \
    | tail -1) || true
  if [[ -n "$last_tx_status" ]]; then
    local last_tx_fn
    last_tx_fn=$(echo "$stats_logs" \
      | grep "Updating tx status" \
      | tail -1 \
      | grep -oP 'function\.name: "\K[^"]+') || last_tx_fn=""
    local last_tx_to
    last_tx_to=$(echo "$stats_logs" \
      | grep "Updating tx status" \
      | tail -1 \
      | grep -oP 'tx\.to: Some\(\K0x[0-9a-fA-F]+') || last_tx_to=""
    local tx_color="$YELLOW"
    [[ "$last_tx_status" == "Finalized" ]] && tx_color="$GREEN"
    echo -e "  Last tx:   ${tx_color}${last_tx_status}${NC}  ${last_tx_fn:+fn=$last_tx_fn}  ${last_tx_to:+to=${last_tx_to:0:14}...}"
  fi
  echo
}

# ─── Warp Token State ───────────────────────────────────────────────────────

print_warp_state() {
  echo -e "${BOLD}Warp Token State${NC}"

  local evm_hyp_erc20 token_id
  evm_hyp_erc20=$(load_state "evm_hyp_erc20")
  token_id=$(load_state "token_id")

  if [[ -n "$evm_hyp_erc20" ]] && command -v cast >/dev/null 2>&1; then
    local supply ism
    supply=$(cast call "$evm_hyp_erc20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null) || supply="?"
    ism=$(cast call "$evm_hyp_erc20" "interchainSecurityModule()(address)" --rpc-url "$EVM_RPC" 2>/dev/null) || ism="?"
    echo -e "  EVM HypERC20:    $evm_hyp_erc20"
    echo -e "  EVM totalSupply: $supply"
    echo -e "  EVM ISM:         $ism"
  else
    echo -e "  EVM HypERC20:    ${YELLOW}not deployed${NC}"
  fi

  if [[ -n "$token_id" ]] && "$BINARY" status --node "$NODE" --home "$BITSONG_HOME" >/dev/null 2>&1; then
    echo -e "  Cosmos token:    $token_id"
  fi
  echo
}

# ─── Test Results ────────────────────────────────────────────────────────────

print_test_results() {
  echo -e "${BOLD}Transfer Tests${NC}"
  local c2e e2c
  c2e=$(load_state "cosmos_to_evm_test_passed")
  e2c=$(load_state "evm_to_cosmos_test_passed")
  if [[ -n "$c2e" ]]; then
    echo -e "  Cosmos -> EVM:  ${GREEN}$c2e${NC}"
  else
    echo -e "  Cosmos -> EVM:  ${YELLOW}not run${NC}"
  fi
  if [[ -n "$e2c" ]]; then
    echo -e "  EVM -> Cosmos:  ${GREEN}$e2c${NC}"
  else
    echo -e "  EVM -> Cosmos:  ${YELLOW}not run${NC}"
  fi
  echo
}

# ─── Show Logs ───────────────────────────────────────────────────────────────

show_logs() {
  local lines="${1:-30}"
  log_step "Relayer Logs (last $lines lines)"

  if ! docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^hyperlane-relayer$"; then
    log_warn "Relayer not running"
    return
  fi

  docker logs hyperlane-relayer --tail "$lines" 2>&1 \
    | grep -v "0xa2827cb39\|CCIP Read" || true
}

# ─── Full Dashboard ─────────────────────────────────────────────────────────

print_dashboard() {
  clear 2>/dev/null || true
  echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════${NC}"
  echo -e "${BOLD}${CYAN}  Hyperlane Status — $(date '+%H:%M:%S')${NC}"
  echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════${NC}"
  echo

  print_containers
  print_chain_heights
  print_checkpoints
  print_relayer_sync
  print_warp_state
  print_test_results
}

# ─── Main ────────────────────────────────────────────────────────────────────

WATCH=false
WATCH_INTERVAL=10
SHOW_LOGS=false
LOG_LINES=30

while [[ $# -gt 0 ]]; do
  case "$1" in
    --watch)
      WATCH=true
      if [[ -n "${2:-}" && "$2" =~ ^[0-9]+$ ]]; then
        WATCH_INTERVAL="$2"; shift
      fi
      shift ;;
    --logs)
      SHOW_LOGS=true
      if [[ -n "${2:-}" && "$2" =~ ^[0-9]+$ ]]; then
        LOG_LINES="$2"; shift
      fi
      shift ;;
    -h|--help)
      echo "Usage: $0 [--watch [interval]] [--logs [lines]]"
      echo ""
      echo "  --watch [N]   Refresh every N seconds (default: 10)"
      echo "  --logs [N]    Show last N lines of relayer logs (default: 30)"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

if [[ "$SHOW_LOGS" == "true" ]]; then
  show_logs "$LOG_LINES"
  exit 0
fi

if [[ "$WATCH" == "true" ]]; then
  trap 'echo; echo "Stopped."; exit 0' INT
  while true; do
    print_dashboard
    echo -e "${CYAN}Refreshing in ${WATCH_INTERVAL}s... (Ctrl+C to stop)${NC}"
    sleep "$WATCH_INTERVAL"
  done
else
  print_dashboard
fi
