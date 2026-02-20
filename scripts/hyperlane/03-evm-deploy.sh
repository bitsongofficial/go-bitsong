#!/usr/bin/env bash
# =============================================================================
# 03-evm-deploy.sh — Deploy HypERC20 on Base Sepolia + bidirectional enrollment
#
# Deploys via Hyperlane CLI, then uses cast for EVM enrollment + verification.
#
# Usage:
#   bash 03-evm-deploy.sh             # Deploy + enroll
#   bash 03-evm-deploy.sh --clean     # Wipe Phase 4 state and redo
#
# Required environment:
#   HYP_KEY (or EVM_PRIVATE_KEY)  — EVM deployer private key (0x...)
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

HYP_CHAINS_DIR="${HYP_CHAINS_DIR:-$HOME/.hyperlane/chains}"

# ─── Step 1: Register Chain Metadata ─────────────────────────────────────────

register_chain_metadata() {
  log_step "Step 1: Register Chain Metadata"

  local mailbox_id; mailbox_id=$(require_state "mailbox_id" "mailbox_id")

  mkdir -p "$HYP_CHAINS_DIR/bitsong"

  cat > "$HYP_CHAINS_DIR/bitsong/metadata.yaml" << EOF
chainId: $CHAIN_ID
domainId: $DOMAIN_ID
name: bitsong
protocol: cosmos
bech32Prefix: bitsong
slip44: 639
rpcUrls:
  - http: http://localhost:26657
restUrls:
  - http: http://localhost:1317
grpcUrls:
  - http: http://localhost:9090
nativeToken:
  name: BitSong
  symbol: BTSG
  decimals: 6
  denom: $DENOM
blocks:
  confirmations: 1
  estimateBlockTime: 6
isTestnet: true
EOF

  cat > "$HYP_CHAINS_DIR/bitsong/addresses.yaml" << EOF
mailbox: "$mailbox_id"
EOF

  log_ok "Chain metadata: $HYP_CHAINS_DIR/bitsong/"
}

# ─── Step 2: Deploy HypERC20 ────────────────────────────────────────────────

deploy_hyp_erc20() {
  log_step "Step 2: Deploy HypERC20 on Base Sepolia"

  DEPLOYED_EVM=$(load_state "evm_hyp_erc20")
  if [[ -n "$DEPLOYED_EVM" && "$DEPLOYED_EVM" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    log_ok "HypERC20 already deployed: $DEPLOYED_EVM"
    return 0
  fi

  # Generate single-chain warp config
  local warp_config="$BITSONG_HOME/warp-route-basesepolia.yaml"
  cat > "$warp_config" << EOF
basesepolia:
  type: synthetic
  name: "BitSong"
  symbol: "BTSG"
  decimals: 6
  mailbox: "$BASESEPOLIA_MAILBOX"
  owner: "$DEPLOYER_ADDR"
  interchainSecurityModule: "0x0000000000000000000000000000000000000000"
EOF

  log "Warp config:"
  cat "$warp_config"; echo

  export HYP_KEY="$EVM_KEY"
  log "Running: hyperlane warp deploy ..."

  # The CLI attempts Etherscan verification (requires API key we don't have).
  hyperlane warp deploy --config "$warp_config" --yes

  # Extract deployed address from CLI artifacts
  log "Looking for deployment artifacts..."
  local latest_warp
  latest_warp=$(find "$HOME/.hyperlane/deployments/warp_routes" -name "*basesepolia*config.yaml" 2>/dev/null \
    | sort -r | head -1) || true

  if [[ -n "$latest_warp" ]]; then
    log "Artifact: $latest_warp"
    DEPLOYED_EVM=$(grep -oP '(?<=addressOrDenom: ")[^"]+' "$latest_warp" 2>/dev/null | head -1) || true
  fi

  if [[ -z "$DEPLOYED_EVM" || ! "$DEPLOYED_EVM" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    log_err "Could not extract HypERC20 address from CLI artifacts."
    log_err "Set manually: jq '.evm_hyp_erc20=\"0xADDR\"' $STATE_FILE > /tmp/s && mv /tmp/s $STATE_FILE"
    exit 1
  fi

  save_state "evm_hyp_erc20" "$DEPLOYED_EVM"
  log_ok "HypERC20: $DEPLOYED_EVM"

  # Compute bytes32
  local evm_bytes32; evm_bytes32=$(evm_to_bytes32 "$DEPLOYED_EVM")
  save_state "evm_hyp_erc20_bytes32" "$evm_bytes32"
  log "  bytes32: $evm_bytes32"
}

# ─── Step 3: Create Cosmos Collateral Token ──────────────────────────────────

create_cosmos_token() {
  log_step "Step 3: Create Cosmos Collateral Token"

  TOKEN_ID=$(load_state "token_id")
  if [[ -n "$TOKEN_ID" ]]; then
    log_ok "Cosmos token already exists: $TOKEN_ID"
    return 0
  fi

  local mailbox_id; mailbox_id=$(require_state "mailbox_id" "mailbox_id")

  submit_tx "MsgCreateCollateralToken (denom=$DENOM)" \
    "$BINARY" tx warp create-collateral-token "$mailbox_id" "$DENOM"

  TOKEN_ID=$(extract_id)
  [[ -n "$TOKEN_ID" ]] || { log_err "Could not extract token_id"; exit 1; }
  save_state "token_id" "$TOKEN_ID"
  log_ok "Cosmos token: $TOKEN_ID"
}

# ─── Step 4: Enroll EVM Router on Cosmos ─────────────────────────────────────

enroll_cosmos_side() {
  log_step "Step 4: Enroll EVM Router on Cosmos"

  local evm_bytes32; evm_bytes32=$(require_state "evm_hyp_erc20_bytes32" "evm_hyp_erc20_bytes32")

  submit_tx "MsgEnrollRemoteRouter (domain=$REMOTE_DOMAIN)" \
    "$BINARY" tx warp enroll-remote-router \
    "$TOKEN_ID" "$REMOTE_DOMAIN" "$evm_bytes32" "$ENROLL_GAS"

  save_state "router_enrolled" "true"
  log_ok "Cosmos side enrolled: domain $REMOTE_DOMAIN -> $evm_bytes32"
}

# ─── Step 5: Enroll Cosmos Router on EVM ─────────────────────────────────────

enroll_evm_side() {
  log_step "Step 5: Enroll Cosmos Router on EVM"

  log "enrollRemoteRouter($DOMAIN_ID, $TOKEN_ID) on $DEPLOYED_EVM ..."
  local tx; tx=$(cast send "$DEPLOYED_EVM" \
    "enrollRemoteRouter(uint32,bytes32)" "$DOMAIN_ID" "$TOKEN_ID" \
    --private-key "$EVM_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true

  local status; status=$(echo "$tx" | jq -r '.status // empty' 2>/dev/null) || true
  local hash;   hash=$(echo "$tx" | jq -r '.transactionHash // empty' 2>/dev/null) || true

  if [[ "$status" == "0x1" || "$status" == "1" ]]; then
    save_state "evm_router_enrolled" "true"
    log_ok "EVM enrolled: domain $DOMAIN_ID -> $TOKEN_ID (tx: $hash)"
  else
    log_err "EVM enrollment failed (status: $status)"
    echo "$tx" | jq -c '.' 2>/dev/null || echo "$tx"
    exit 1
  fi
}

# ─── Step 6: Verify ─────────────────────────────────────────────────────────

verify_deployment() {
  log_step "Step 6: Verify"

  log "--- EVM HypERC20 ($DEPLOYED_EVM) ---"
  local decimals mailbox supply router_on_evm
  decimals=$(cast call "$DEPLOYED_EVM" "decimals()(uint8)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  mailbox=$(cast call "$DEPLOYED_EVM" "mailbox()(address)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  supply=$(cast call "$DEPLOYED_EVM" "totalSupply()(uint256)" --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")
  router_on_evm=$(cast call "$DEPLOYED_EVM" "routers(uint32)(bytes32)" "$DOMAIN_ID" \
    --rpc-url "$EVM_RPC" 2>/dev/null || echo "?")

  echo "  decimals:          $decimals  (expected 6)"
  echo "  mailbox:           $mailbox"
  echo "  totalSupply:       $supply   (expected 0)"
  echo "  router($DOMAIN_ID): $router_on_evm"
  echo

  log "--- Cosmos warp token ($TOKEN_ID) ---"
  "$BINARY" query warp remote-routers "$TOKEN_ID" \
    --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | jq '.' || true
}

# ─── Summary ─────────────────────────────────────────────────────────────────

print_summary() {
  log_step "Summary"
  echo -e "${BOLD}Warp Route${NC}"
  echo "  Cosmos token:  $(load_state token_id)"
  echo "  EVM HypERC20:  $(load_state evm_hyp_erc20)"
  echo "  EVM (bytes32): $(load_state evm_hyp_erc20_bytes32)"
  echo "  State file:    $STATE_FILE"
  echo
  echo -e "${BOLD}BaseScan${NC}: https://sepolia.basescan.org/address/$(load_state evm_hyp_erc20)"
  echo
  echo -e "${BOLD}Next${NC}: Run 04-agents.sh to start validators + relayer"
}

# ─── Main ────────────────────────────────────────────────────────────────────

CLEAN=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --clean) CLEAN=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--clean]"
      echo "  --clean   Wipe Phase 4 state keys and redo"
      echo ""
      echo "Required: HYP_KEY or EVM_PRIVATE_KEY (deployer private key)"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

# Preflight
require_jq
command -v hyperlane >/dev/null 2>&1 || { log_err "Hyperlane CLI missing: npm install -g @hyperlane-xyz/cli"; exit 1; }
command -v cast >/dev/null 2>&1 || { log_err "cast missing: foundryup"; exit 1; }

[[ -f "$STATE_FILE" ]] || { log_err "Phase 3 state not found: $STATE_FILE. Run 02-hyperlane.sh first."; exit 1; }

EVM_KEY="${HYP_KEY:-${EVM_PRIVATE_KEY:-}}"
[[ -n "$EVM_KEY" ]] || { log_err "Set HYP_KEY=0x<private_key_hex>"; exit 1; }

DEPLOYER_ADDR=$(cast wallet address --private-key "$EVM_KEY" 2>/dev/null) \
  || { log_err "Cannot derive EVM address from key"; exit 1; }

require_binary
require_chain_running

log_ok "Hyperlane CLI: $(hyperlane --version 2>&1 | head -1)"
log_ok "cast: $(cast --version 2>&1 | head -1)"
log_ok "EVM deployer: $DEPLOYER_ADDR"
log_ok "Phase 3 mailbox: $(load_state mailbox_id)"

if [[ "$CLEAN" == "true" ]]; then
  jq 'del(.token_id, .router_enrolled, .evm_hyp_erc20, .evm_hyp_erc20_bytes32, .evm_router_enrolled)' \
    "$STATE_FILE" > "${STATE_FILE}.tmp" && mv "${STATE_FILE}.tmp" "$STATE_FILE"
  log_ok "Phase 4 state cleared"
fi

banner "Phase 4: Deploy Warp Route" "bitsong <-> Base Sepolia"

register_chain_metadata
deploy_hyp_erc20
create_cosmos_token
enroll_cosmos_side
enroll_evm_side
verify_deployment
print_summary

log_ok "Phase 4 complete!"
