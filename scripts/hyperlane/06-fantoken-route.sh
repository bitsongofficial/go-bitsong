#!/usr/bin/env bash
# =============================================================================
# 06-fantoken-route.sh — Issue a fantoken and create its Hyperlane warp route
#
# Supports multiple fantoken warp routes, each stored independently in state
# using per-symbol prefixed keys (ft_<sym>_*).
#
# Usage:
#   bash 06-fantoken-route.sh --symbol clay --name "Clay Token"   # Register new
#   bash 06-fantoken-route.sh --symbol clay                       # Resume/show
#   bash 06-fantoken-route.sh --list                              # List all routes
#   bash 06-fantoken-route.sh --symbol clay --clean               # Wipe & redo
#
# Required environment:
#   HYP_KEY (or EVM_PRIVATE_KEY)  — EVM deployer private key (0x...)
#   EVM_RELAYER_KEY               — for relayer restart
#   VALIDATOR_KEY, COSMOS_SIGNER_KEY — for relayer restart
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

HYP_KEY="${HYP_KEY:-${EVM_PRIVATE_KEY:-}}"
EVM_RELAYER_KEY="${EVM_RELAYER_KEY:-${HYP_KEY:-}}"

FT_MAX_SUPPLY="${FT_MAX_SUPPLY:-1000000000000}"
FT_MINT_AMOUNT="${FT_MINT_AMOUNT:-1000000000}"

HYP_CHAINS_DIR="${HYP_CHAINS_DIR:-$HOME/.hyperlane/chains}"

# =============================================================================
# EVM nonce helper — wait for pending TX to confirm
# =============================================================================

# Wait until EVM pending nonce equals confirmed nonce (no pending TXs).
# Prevents "replacement transaction underpriced" when sending back-to-back.
wait_evm_pending() {
  local addr
  addr=$(cast wallet address --private-key "$HYP_KEY" 2>/dev/null) || return 0
  log "Waiting for pending EVM TXs to confirm ($addr)..."
  for _ in $(seq 1 30); do
    local pending confirmed
    pending=$(cast nonce "$addr" --rpc-url "$EVM_RPC" --pending 2>/dev/null) || break
    confirmed=$(cast nonce "$addr" --rpc-url "$EVM_RPC" 2>/dev/null) || break
    if [[ "$pending" == "$confirmed" ]]; then
      log_ok "EVM nonce settled (nonce=$confirmed)"
      return 0
    fi
    sleep 3
  done
  log_warn "Timed out waiting for EVM pending TXs (proceeding anyway)"
}

# =============================================================================
# Per-symbol state helpers
# =============================================================================

FT_KEY=""  # lowercased symbol, set by CLI parsing

ft_save() { save_state "ft_${FT_KEY}_$1" "$2"; }
ft_load() { load_state "ft_${FT_KEY}_$1"; }

# Write flat aliases so test/status scripts work with the latest route
ft_set_current() {
  save_state "ft_denom" "$(ft_load denom)"
  save_state "ft_minted" "$(ft_load minted)"
  save_state "ft_token_id" "$(ft_load token_id)"
  save_state "ft_evm_hyp_erc20" "$(ft_load evm_hyp_erc20)"
  save_state "ft_evm_hyp_erc20_bytes32" "$(ft_load evm_hyp_erc20_bytes32)"
  save_state "ft_router_enrolled" "$(ft_load router_enrolled)"
  save_state "ft_evm_router_enrolled" "$(ft_load evm_router_enrolled)"
  save_state "ft_evm_ism_set" "$(ft_load evm_ism_set)"
  save_state "ft_relayer_restarted" "$(ft_load relayer_restarted)"
}

# Add symbol to ft_route_list if not already present
register_route() {
  local list
  list=$(load_state "ft_route_list")
  if [[ -z "$list" ]]; then
    save_state "ft_route_list" "$FT_KEY"
  elif ! echo ",$list," | grep -q ",$FT_KEY,"; then
    save_state "ft_route_list" "${list},${FT_KEY}"
  fi
}

# Remove symbol from ft_route_list
unregister_route() {
  local list
  list=$(load_state "ft_route_list")
  [[ -z "$list" ]] && return 0
  # Remove the symbol and clean up commas
  local new_list
  new_list=$(echo "$list" | tr ',' '\n' | grep -v "^${FT_KEY}$" | paste -sd ',' -) || true
  save_state "ft_route_list" "$new_list"
}

# =============================================================================
# Migration: detect old flat ft_* keys with no ft_route_list
# =============================================================================

migrate_legacy_state() {
  local route_list
  route_list=$(load_state "ft_route_list")
  [[ -n "$route_list" ]] && return 0  # Already migrated

  local old_denom
  old_denom=$(load_state "ft_denom")
  [[ -z "$old_denom" ]] && return 0  # No legacy state

  log "Migrating legacy fantoken state..."

  # Try to derive symbol from chain query
  local old_symbol=""
  if "$BINARY" status --node "$NODE" --home "$BITSONG_HOME" >/dev/null 2>&1; then
    old_symbol=$("$BINARY" query fantoken denom "$old_denom" \
      --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
      jq -r '.fantoken.symbol // empty' 2>/dev/null) || true
  fi
  if [[ -z "$old_symbol" ]]; then
    old_symbol="legacy"
    log_warn "Could not query symbol for $old_denom, using 'legacy'"
  fi
  old_symbol=$(echo "$old_symbol" | tr '[:upper:]' '[:lower:]')

  # Copy flat keys to per-symbol keys
  for suffix in denom minted token_id evm_hyp_erc20 evm_hyp_erc20_bytes32 \
                router_enrolled evm_router_enrolled evm_ism_set relayer_restarted; do
    local val
    val=$(load_state "ft_${suffix}")
    [[ -n "$val" ]] && save_state "ft_${old_symbol}_${suffix}" "$val"
  done

  save_state "ft_route_list" "$old_symbol"
  log_ok "Migrated legacy state to symbol '$old_symbol'"
}

# =============================================================================
# --list: show all registered routes
# =============================================================================

list_routes() {
  local route_list
  route_list=$(load_state "ft_route_list")
  if [[ -z "$route_list" ]]; then
    echo "No fantoken routes registered."
    return 0
  fi

  echo -e "${BOLD}Registered Fantoken Routes${NC}"
  echo -e "  SYMBOL     DENOM                                          EVM HYERC20                                  STATUS"
  echo -e "  ──────── ────────────────────────────────────────────── ────────────────────────────────────────────── ──────────"

  IFS=',' read -ra routes <<< "$route_list"
  for sym in "${routes[@]}"; do
    [[ -z "$sym" ]] && continue
    local denom evm_hyp status_text
    denom=$(load_state "ft_${sym}_denom")
    evm_hyp=$(load_state "ft_${sym}_evm_hyp_erc20")
    local ism_set
    ism_set=$(load_state "ft_${sym}_evm_ism_set")

    if [[ -n "$evm_hyp" && "$ism_set" == "true" ]]; then
      status_text="${GREEN}complete${NC}"
    elif [[ -n "$evm_hyp" ]]; then
      status_text="${YELLOW}partial: no ISM${NC}"
    elif [[ -n "$denom" ]]; then
      status_text="${YELLOW}partial: no EVM deploy${NC}"
    else
      status_text="${RED}empty${NC}"
    fi

    printf "  %-10s %-50s %-44s " "$sym" "${denom:-—}" "${evm_hyp:-—}"
    echo -e "$status_text"
  done
  echo
}

# =============================================================================
# Build dynamic relayer whitelist from all routes
# =============================================================================

build_whitelist() {
  local whitelist="["

  # Always include ubtsg route
  local token_id evm_hyp_erc20
  token_id=$(load_state "token_id")
  evm_hyp_erc20=$(load_state "evm_hyp_erc20")
  if [[ -n "$token_id" && -n "$evm_hyp_erc20" ]]; then
    whitelist+="{\"senderAddress\":\"${token_id}\",\"destinationDomain\":\"${REMOTE_DOMAIN}\"},"
    whitelist+="{\"senderAddress\":\"$(evm_to_bytes32 "$evm_hyp_erc20")\",\"destinationDomain\":\"${DOMAIN_ID}\"},"
  fi

  # Add each fantoken route
  local route_list
  route_list=$(load_state "ft_route_list")
  if [[ -n "$route_list" ]]; then
    IFS=',' read -ra routes <<< "$route_list"
    for sym in "${routes[@]}"; do
      [[ -z "$sym" ]] && continue
      local tid evm
      tid=$(load_state "ft_${sym}_token_id")
      evm=$(load_state "ft_${sym}_evm_hyp_erc20")
      [[ -n "$tid" && -n "$evm" ]] || continue
      whitelist+="{\"senderAddress\":\"${tid}\",\"destinationDomain\":\"${REMOTE_DOMAIN}\"},"
      whitelist+="{\"senderAddress\":\"$(evm_to_bytes32 "$evm")\",\"destinationDomain\":\"${DOMAIN_ID}\"},"
    done
  fi

  whitelist="${whitelist%,}]"  # Remove trailing comma
  echo "$whitelist"
}

# =============================================================================
# Step 1: Issue Fantoken
# =============================================================================

issue_fantoken() {
  log_step "Step 1: Issue Fantoken ($FT_SYMBOL)"

  FT_DENOM=$(ft_load "denom")
  if [[ -n "$FT_DENOM" ]]; then
    log_ok "Fantoken already issued: $FT_DENOM"
    return 0
  fi

  [[ -n "$FT_NAME" ]] || { log_err "--name is required when issuing a new fantoken"; return 1; }

  submit_tx "Issue fantoken ($FT_SYMBOL)" \
    "$BINARY" tx fantoken issue \
    --symbol "$FT_SYMBOL" --name "$FT_NAME" --max-supply "$FT_MAX_SUPPLY"

  # Extract denom from protobuf data field (MsgIssueResponse contains the denom)
  local data
  data=$(echo "$TX_RESULT" | jq -r '.data // empty' 2>/dev/null) || true
  if [[ -n "$data" ]]; then
    FT_DENOM=$(echo "$data" | xxd -r -p | grep -aoP 'ft[0-9A-Fa-f]{40}' | head -1) || true
  fi

  # Fallback: query fantokens by authority
  if [[ -z "$FT_DENOM" ]]; then
    log "Extracting denom from data failed, querying by authority..."
    FT_DENOM=$("$BINARY" query fantoken authority "$VAL_ADDRESS" \
      --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
      jq -r --arg sym "$FT_SYMBOL" '
        [.fantokens[]? | select(.symbol == $sym) | .denom] | first // empty
      ' 2>/dev/null) || true
  fi

  [[ -n "$FT_DENOM" ]] || { log_err "Could not extract fantoken denom"; return 1; }
  ft_save "denom" "$FT_DENOM"
  log_ok "Fantoken issued: $FT_DENOM"
}

# =============================================================================
# Step 2: Mint Fantoken
# =============================================================================

mint_fantoken() {
  log_step "Step 2: Mint Fantoken to Validator ($FT_SYMBOL)"

  if [[ "$(ft_load "minted")" == "true" ]]; then
    log_ok "Already minted"
    return 0
  fi

  submit_tx "Mint ${FT_MINT_AMOUNT} ${FT_DENOM}" \
    "$BINARY" tx fantoken mint "${FT_MINT_AMOUNT}${FT_DENOM}" \
    --recipient "$VAL_ADDRESS"

  ft_save "minted" "true"

  local balance
  balance=$("$BINARY" query bank balance "$VAL_ADDRESS" "$FT_DENOM" \
    --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
    jq -r '.balance.amount // "?"' 2>/dev/null) || balance="?"
  log_ok "Minted. Validator balance: ${balance} ${FT_DENOM}"
}

# =============================================================================
# Step 3: Create Collateral Token
# =============================================================================

create_ft_collateral() {
  log_step "Step 3: Create Collateral Token ($FT_SYMBOL)"

  FT_TOKEN_ID=$(ft_load "token_id")
  if [[ -n "$FT_TOKEN_ID" ]]; then
    log_ok "Collateral token exists: $FT_TOKEN_ID"
    return 0
  fi

  local mailbox_id; mailbox_id=$(require_state "mailbox_id" "mailbox_id")

  submit_tx "Create collateral token ($FT_DENOM)" \
    "$BINARY" tx warp create-collateral-token "$mailbox_id" "$FT_DENOM"

  FT_TOKEN_ID=$(extract_id)
  [[ -n "$FT_TOKEN_ID" ]] || { log_err "Could not extract ft_token_id"; return 1; }
  ft_save "token_id" "$FT_TOKEN_ID"
  log_ok "Cosmos token: $FT_TOKEN_ID"
}

# =============================================================================
# Step 4: Deploy HypERC20 on Base Sepolia
# =============================================================================

deploy_ft_hyp_erc20() {
  log_step "Step 4: Deploy HypERC20 for $FT_SYMBOL"

  FT_EVM_HYP_ERC20=$(ft_load "evm_hyp_erc20")
  if [[ -n "$FT_EVM_HYP_ERC20" && "$FT_EVM_HYP_ERC20" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    log_ok "HypERC20 already deployed: $FT_EVM_HYP_ERC20"
    return 0
  fi

  local deployer_addr
  deployer_addr=$(cast wallet address --private-key "$HYP_KEY" 2>/dev/null) \
    || { log_err "Cannot derive EVM address from HYP_KEY"; return 1; }

  # Generate single-chain warp config for fantoken
  local warp_config="$BITSONG_HOME/warp-route-ft-${FT_KEY}-basesepolia.yaml"
  cat > "$warp_config" << EOF
basesepolia:
  type: synthetic
  name: "$FT_NAME"
  symbol: "$FT_SYMBOL"
  decimals: 6
  mailbox: "$BASESEPOLIA_MAILBOX"
  owner: "$deployer_addr"
  interchainSecurityModule: "0x0000000000000000000000000000000000000000"
EOF

  log "Warp config:"
  cat "$warp_config"; echo

  export HYP_KEY
  log "Running: hyperlane warp deploy ..."
  hyperlane warp deploy --config "$warp_config" --yes

  # Extract deployed address from CLI artifacts
  log "Looking for deployment artifacts..."
  local latest_warp
  latest_warp=$(find "$HOME/.hyperlane/deployments/warp_routes" -name "*basesepolia*config.yaml" 2>/dev/null \
    | sort -r | head -1) || true

  if [[ -n "$latest_warp" ]]; then
    log "Artifact: $latest_warp"
    FT_EVM_HYP_ERC20=$(grep -oP '(?<=addressOrDenom: ")[^"]+' "$latest_warp" 2>/dev/null | head -1) || true
  fi

  if [[ -z "$FT_EVM_HYP_ERC20" || ! "$FT_EVM_HYP_ERC20" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    log_err "Could not extract HypERC20 address from CLI artifacts."
    log_err "Set manually: jq '.ft_${FT_KEY}_evm_hyp_erc20=\"0xADDR\"' $STATE_FILE > /tmp/s && mv /tmp/s $STATE_FILE"
    return 1
  fi

  ft_save "evm_hyp_erc20" "$FT_EVM_HYP_ERC20"
  log_ok "FT HypERC20: $FT_EVM_HYP_ERC20"

  local evm_bytes32; evm_bytes32=$(evm_to_bytes32 "$FT_EVM_HYP_ERC20")
  ft_save "evm_hyp_erc20_bytes32" "$evm_bytes32"
  log "  bytes32: $evm_bytes32"
}

# =============================================================================
# Step 5: Enroll Remote Routers (bidirectional)
# =============================================================================

enroll_ft_routers() {
  log_step "Step 5: Enroll Remote Routers ($FT_SYMBOL)"

  # Cosmos side: enroll EVM HypERC20 as remote router
  if [[ "$(ft_load "router_enrolled")" != "true" ]]; then
    local evm_bytes32
    evm_bytes32=$(ft_load "evm_hyp_erc20_bytes32")
    [[ -n "$evm_bytes32" ]] || { log_err "ft_${FT_KEY}_evm_hyp_erc20_bytes32 not in state"; return 1; }

    submit_tx "Enroll FT remote router on Cosmos (domain=$REMOTE_DOMAIN)" \
      "$BINARY" tx warp enroll-remote-router \
      "$FT_TOKEN_ID" "$REMOTE_DOMAIN" "$evm_bytes32" "$ENROLL_GAS"

    ft_save "router_enrolled" "true"
    log_ok "Cosmos enrolled: domain $REMOTE_DOMAIN -> $evm_bytes32"
  else
    log_ok "Cosmos router already enrolled"
  fi

  # EVM side: enroll Cosmos token as remote router
  if [[ "$(ft_load "evm_router_enrolled")" != "true" ]]; then
    log "enrollRemoteRouter($DOMAIN_ID, $FT_TOKEN_ID) on $FT_EVM_HYP_ERC20 ..."
    local tx; tx=$(cast send "$FT_EVM_HYP_ERC20" \
      "enrollRemoteRouter(uint32,bytes32)" "$DOMAIN_ID" "$FT_TOKEN_ID" \
      --private-key "$HYP_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true

    local status; status=$(echo "$tx" | jq -r '.status // empty' 2>/dev/null) || true
    local hash;   hash=$(echo "$tx" | jq -r '.transactionHash // empty' 2>/dev/null) || true

    if [[ "$status" == "0x1" || "$status" == "1" ]]; then
      ft_save "evm_router_enrolled" "true"
      log_ok "EVM enrolled: domain $DOMAIN_ID -> $FT_TOKEN_ID (tx: $hash)"
    else
      log_err "EVM enrollment failed (status: $status)"
      echo "$tx" | jq -c '.' 2>/dev/null || echo "$tx"
      return 1
    fi
  else
    log_ok "EVM router already enrolled"
  fi
}

# =============================================================================
# Step 6: Set ISM on EVM HypERC20
# =============================================================================

set_ft_ism() {
  log_step "Step 6: Set ISM on FT HypERC20 ($FT_SYMBOL)"

  if [[ "$(ft_load "evm_ism_set")" == "true" ]]; then
    log_ok "ISM already set"
    return 0
  fi

  local basesepolia_multisig_ism
  basesepolia_multisig_ism=$(require_state "basesepolia_multisig_ism" "basesepolia_multisig_ism (run 04-agents.sh first)")

  log "Setting ISM to $basesepolia_multisig_ism on $FT_EVM_HYP_ERC20..."
  local tx; tx=$(cast send "$FT_EVM_HYP_ERC20" \
    "setInterchainSecurityModule(address)" "$basesepolia_multisig_ism" \
    --private-key "$HYP_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true

  local status; status=$(echo "$tx" | jq -r '.status // empty' 2>/dev/null) || true
  if [[ "$status" == "0x1" || "$status" == "1" ]]; then
    ft_save "evm_ism_set" "true"
    log_ok "ISM set to $basesepolia_multisig_ism"
  else
    log_err "setInterchainSecurityModule failed (status: $status)"
    echo "$tx" | jq -c '.' 2>/dev/null || echo "$tx"
    return 1
  fi
}

# =============================================================================
# Step 7: Restart Relayer with Dynamic Whitelist
# =============================================================================

restart_relayer() {
  log_step "Step 7: Restart Relayer (dynamic whitelist)"

  # Always restart when a new route is added — the whitelist changes
  local whitelist
  whitelist=$(build_whitelist)

  # Count entries
  local entry_count
  entry_count=$(echo "$whitelist" | jq 'length' 2>/dev/null) || entry_count="?"

  log "Stopping existing relayer..."
  docker stop hyperlane-relayer 2>/dev/null || true
  docker rm -f hyperlane-relayer 2>/dev/null || true
  sleep 2

  log "Starting relayer with dynamic whitelist ($entry_count entries)..."
  log "  Whitelist: $whitelist"
  docker run -d --name hyperlane-relayer --network host \
    -e CONFIG_FILES=/config/agent-config.json \
    -v "$BITSONG_HOME/agent-config.json:/config/agent-config.json:ro" \
    -v "$BITSONG_HOME/relayer-db:/hyperlane_db" \
    -v "$BITSONG_HOME/checkpoints-bitsong:/checkpoints-bitsong:ro" \
    -v "$BITSONG_HOME/checkpoints-basesepolia:/checkpoints-basesepolia:ro" \
    "$DOCKER_IMAGE" ./relayer \
    --db /hyperlane_db --relayChains bitsong,basesepolia \
    --allowLocalCheckpointSyncers true \
    --gaspaymentenforcement '[{"type": "none"}]' \
    --whitelist "$whitelist" \
    --chains.bitsong.signer.type cosmosKey \
    --chains.bitsong.signer.key "$COSMOS_SIGNER_KEY" \
    --chains.bitsong.signer.prefix bitsong \
    --chains.basesepolia.signer.type hexKey \
    --chains.basesepolia.signer.key "$EVM_RELAYER_KEY" \
    --metricsPort 9091

  sleep 3
  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^hyperlane-relayer$"; then
    ft_save "relayer_restarted" "true"
    log_ok "Relayer restarted with $entry_count whitelist entries"
  else
    log_err "Relayer failed to start — check: docker logs hyperlane-relayer"
    return 1
  fi
}

# =============================================================================
# Verify
# =============================================================================

verify_ft_deployment() {
  log_step "Verify Fantoken Warp Route ($FT_SYMBOL)"

  log "--- Cosmos Warp Token ($FT_TOKEN_ID) ---"
  "$BINARY" query warp remote-routers "$FT_TOKEN_ID" \
    --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | jq '.' || true
  echo

  log "--- EVM HypERC20 ($FT_EVM_HYP_ERC20) ---"
  local decimals symbol mailbox supply ism router_on_evm
  symbol=$(cast call "$FT_EVM_HYP_ERC20" "symbol()(string)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  decimals=$(cast call "$FT_EVM_HYP_ERC20" "decimals()(uint8)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  mailbox=$(cast call "$FT_EVM_HYP_ERC20" "mailbox()(address)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  supply=$(cast call "$FT_EVM_HYP_ERC20" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  ism=$(cast call "$FT_EVM_HYP_ERC20" "interchainSecurityModule()(address)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  router_on_evm=$(cast call "$FT_EVM_HYP_ERC20" "routers(uint32)(bytes32)" "$DOMAIN_ID" \
    --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")

  echo "  symbol:            $symbol  (expected $FT_SYMBOL)"
  echo "  decimals:          $decimals  (expected 6)"
  echo "  mailbox:           $mailbox"
  echo "  totalSupply:       $supply  (expected 0)"
  echo "  ISM:               $ism"
  echo "  router($DOMAIN_ID): $router_on_evm"
  echo
}

# =============================================================================
# Summary
# =============================================================================

print_summary() {
  log_step "Summary"
  echo -e "${BOLD}Fantoken Warp Route ($FT_SYMBOL)${NC}"
  echo "  Denom:         $(ft_load denom)"
  echo "  Cosmos token:  $(ft_load token_id)"
  echo "  EVM HypERC20:  $(ft_load evm_hyp_erc20)"
  echo "  EVM (bytes32): $(ft_load evm_hyp_erc20_bytes32)"
  echo "  State prefix:  ft_${FT_KEY}_*"
  echo "  State file:    $STATE_FILE"
  echo
  echo -e "${BOLD}BaseScan${NC}: https://sepolia.basescan.org/address/$(ft_load evm_hyp_erc20)"
  echo
  echo -e "${BOLD}Next${NC}: Run 07-fantoken-test.sh --symbol $FT_SYMBOL to test transfers"
}

# =============================================================================
# Main
# =============================================================================

CLEAN=false
LIST_MODE=false
FT_SYMBOL=""
FT_NAME=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --symbol)
      [[ -n "${2:-}" ]] || { log_err "--symbol requires a value"; exit 1; }
      FT_SYMBOL="$2"; shift 2 ;;
    --name)
      [[ -n "${2:-}" ]] || { log_err "--name requires a value"; exit 1; }
      FT_NAME="$2"; shift 2 ;;
    --list) LIST_MODE=true; shift ;;
    --clean) CLEAN=true; shift ;;
    -h|--help)
      echo "Usage: $0 --symbol <sym> [--name \"Name\"] [--clean]"
      echo "       $0 --list"
      echo ""
      echo "Options:"
      echo "  --symbol <sym>   Fantoken symbol (required unless --list)"
      echo "  --name \"Name\"    Fantoken name (required for new tokens)"
      echo "  --list           List all registered fantoken routes"
      echo "  --clean          Wipe this symbol's state keys and redo"
      echo ""
      echo "Required env: HYP_KEY, EVM_RELAYER_KEY, VALIDATOR_KEY, COSMOS_SIGNER_KEY"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

# Handle --list early (only needs state file)
if [[ "$LIST_MODE" == "true" ]]; then
  [[ -f "$STATE_FILE" ]] || { echo "No state file found."; exit 0; }
  migrate_legacy_state
  list_routes
  exit 0
fi

# --symbol is required for everything else
[[ -n "$FT_SYMBOL" ]] || { log_err "--symbol is required (or use --list)"; exit 1; }
FT_KEY=$(echo "$FT_SYMBOL" | tr '[:upper:]' '[:lower:]')

# Preflight
require_binary
require_jq
require_chain_running
command -v hyperlane >/dev/null 2>&1 || { log_err "Hyperlane CLI missing: npm install -g @hyperlane-xyz/cli"; exit 1; }
command -v cast >/dev/null 2>&1 || { log_err "cast missing: foundryup"; exit 1; }

[[ -f "$STATE_FILE" ]] || { log_err "State not found: $STATE_FILE. Run earlier phases first."; exit 1; }
[[ -n "$HYP_KEY" ]] || { log_err "Set HYP_KEY=0x<private_key_hex>"; exit 1; }
[[ -n "${VALIDATOR_KEY:-}" ]]     || { log_err "VALIDATOR_KEY not set (needed for relayer restart)"; exit 1; }
[[ -n "${COSMOS_SIGNER_KEY:-}" ]] || { log_err "COSMOS_SIGNER_KEY not set (needed for relayer restart)"; exit 1; }
[[ -n "${EVM_RELAYER_KEY:-}" ]]   || { log_err "EVM_RELAYER_KEY not set (needed for relayer restart)"; exit 1; }
[[ -n "$VAL_ADDRESS" ]] || { log_err "VAL_ADDRESS not available. Is chain initialized?"; exit 1; }

# Migrate legacy flat keys if needed
migrate_legacy_state

# Handle --clean for this specific symbol
if [[ "$CLEAN" == "true" ]]; then
  log "Cleaning state for symbol '$FT_KEY'..."
  for suffix in denom minted token_id evm_hyp_erc20 evm_hyp_erc20_bytes32 \
               router_enrolled evm_router_enrolled evm_ism_set relayer_restarted; do
    jq --arg k "ft_${FT_KEY}_${suffix}" 'del(.[$k])' "$STATE_FILE" > "${STATE_FILE}.tmp" \
      && mv "${STATE_FILE}.tmp" "$STATE_FILE"
  done
  unregister_route
  log_ok "State cleaned for '$FT_KEY'"
fi

banner "Phase 6: Fantoken Warp Route" "$FT_SYMBOL on bitsong <-> basesepolia"

issue_fantoken
mint_fantoken
create_ft_collateral
deploy_ft_hyp_erc20
wait_evm_pending  # ensure HypERC20 deploy TX confirms before enrollment
enroll_ft_routers
wait_evm_pending  # ensure EVM enrollment TX confirms before ISM TX
set_ft_ism

# Register route BEFORE relayer restart so build_whitelist() includes it
register_route
ft_set_current

restart_relayer
verify_ft_deployment
print_summary

log_ok "Phase 6 complete for $FT_SYMBOL!"
