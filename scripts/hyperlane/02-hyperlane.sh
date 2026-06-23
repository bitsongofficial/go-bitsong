#!/usr/bin/env bash
# =============================================================================
# 02-hyperlane.sh — Configure Hyperlane bridge on the local chain (9 txs)
#
# Usage:
#   bash 02-hyperlane.sh                                    # Full 9-step setup
#   bash 02-hyperlane.sh --update-enrollment 0xABC...DEF    # Re-enroll with real EVM addr
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

# EVM_RECEIVER_CONTRACT is only used if set externally (not enrolled with zero placeholder)
EVM_RECEIVER_CONTRACT="${EVM_RECEIVER_CONTRACT:-}"

# ─── 9-Step Hyperlane Setup ──────────────────────────────────────────────────

setup_hyperlane() {
  log_step "Configuring Hyperlane Bridge (9 Steps)"

  local multisig_ism_id routing_ism_id mailbox_id merkle_hook_id igp_id token_id

  # Step 1: MultisigISM
  multisig_ism_id=$(load_state "multisig_ism_id")
  if [[ -n "$multisig_ism_id" ]]; then
    log_ok "1/9: MultisigISM exists ($multisig_ism_id)"
  else
    submit_tx "1/9: Create MessageIdMultisigISM" \
      "$BINARY" tx hyperlane ism create-message-id-multisig \
      "$VALIDATOR_ADDRESSES" "$VALIDATOR_THRESHOLD"
    multisig_ism_id=$(extract_id)
    [[ -n "$multisig_ism_id" ]] || { log_err "Failed to extract MultisigISM ID"; return 1; }
    save_state "multisig_ism_id" "$multisig_ism_id"
    log_ok "1/9: MultisigISM = $multisig_ism_id"
  fi

  # Step 2: RoutingISM
  routing_ism_id=$(load_state "routing_ism_id")
  if [[ -n "$routing_ism_id" ]]; then
    log_ok "2/9: RoutingISM exists ($routing_ism_id)"
  else
    submit_tx "2/9: Create RoutingISM" \
      "$BINARY" tx hyperlane ism create-routing \
      --routes="[{\"domain\":${REMOTE_DOMAIN},\"ism\":\"${multisig_ism_id}\"}]"
    routing_ism_id=$(extract_id)
    [[ -n "$routing_ism_id" ]] || { log_err "Failed to extract RoutingISM ID"; return 1; }
    save_state "routing_ism_id" "$routing_ism_id"
    log_ok "2/9: RoutingISM = $routing_ism_id"
  fi

  # Step 3: Mailbox
  mailbox_id=$(load_state "mailbox_id")
  if [[ -n "$mailbox_id" ]]; then
    log_ok "3/9: Mailbox exists ($mailbox_id)"
  else
    submit_tx "3/9: Create Mailbox (domain=$DOMAIN_ID)" \
      "$BINARY" tx hyperlane mailbox create "$routing_ism_id" "$DOMAIN_ID"
    mailbox_id=$(extract_id)
    [[ -n "$mailbox_id" ]] || { log_err "Failed to extract Mailbox ID"; return 1; }
    save_state "mailbox_id" "$mailbox_id"
    log_ok "3/9: Mailbox = $mailbox_id"
  fi

  # Step 4: MerkleTreeHook
  merkle_hook_id=$(load_state "merkle_hook_id")
  if [[ -n "$merkle_hook_id" ]]; then
    log_ok "4/9: MerkleTreeHook exists ($merkle_hook_id)"
  else
    submit_tx "4/9: Create MerkleTreeHook" \
      "$BINARY" tx hyperlane hooks merkle create "$mailbox_id"
    merkle_hook_id=$(extract_id)
    [[ -n "$merkle_hook_id" ]] || { log_err "Failed to extract MerkleTreeHook ID"; return 1; }
    save_state "merkle_hook_id" "$merkle_hook_id"
    log_ok "4/9: MerkleTreeHook = $merkle_hook_id"
  fi

  # Step 5: IGP
  igp_id=$(load_state "igp_id")
  if [[ -n "$igp_id" ]]; then
    log_ok "5/9: IGP exists ($igp_id)"
  else
    submit_tx "5/9: Create IGP (denom=$DENOM)" \
      "$BINARY" tx hyperlane hooks igp create "$DENOM"
    igp_id=$(extract_id)
    [[ -n "$igp_id" ]] || { log_err "Failed to extract IGP ID"; return 1; }
    save_state "igp_id" "$igp_id"
    log_ok "5/9: IGP = $igp_id"
  fi

  # Step 6: IGP Gas Oracle
  if [[ "$(load_state "igp_configured")" == "true" ]]; then
    log_ok "6/9: IGP gas config already set"
  else
    submit_tx "6/9: Set IGP destination gas config (remote=$REMOTE_DOMAIN)" \
      "$BINARY" tx hyperlane hooks igp set-destination-gas-config \
      "$igp_id" "$REMOTE_DOMAIN" "$TOKEN_EXCHANGE_RATE" "$GAS_PRICE" "$GAS_OVERHEAD"
    save_state "igp_configured" "true"
    log_ok "6/9: IGP gas config set for domain $REMOTE_DOMAIN"
  fi

  # Step 7: Set Mailbox Hooks
  if [[ "$(load_state "mailbox_updated")" == "true" ]]; then
    log_ok "7/9: Mailbox hooks already set"
  else
    submit_tx "7/9: Set Mailbox hooks (default=IGP, required=MerkleTree)" \
      "$BINARY" tx hyperlane mailbox set "$mailbox_id" \
      --default-hook="$igp_id" --required-hook="$merkle_hook_id"
    save_state "mailbox_updated" "true"
    log_ok "7/9: Mailbox updated with hooks"
  fi

  # Step 8: Collateral Token
  token_id=$(load_state "token_id")
  if [[ -n "$token_id" ]]; then
    log_ok "8/9: Collateral token exists ($token_id)"
  else
    submit_tx "8/9: Create collateral token ($DENOM)" \
      "$BINARY" tx warp create-collateral-token "$mailbox_id" "$DENOM"
    token_id=$(extract_id)
    [[ -n "$token_id" ]] || { log_err "Failed to extract Token ID"; return 1; }
    save_state "token_id" "$token_id"
    log_ok "8/9: Token = $token_id"
  fi

  # Step 9: Enroll Remote Router (only if real EVM address is available)
  if [[ "$(load_state "router_enrolled")" == "true" ]]; then
    log_ok "9/9: Remote router already enrolled"
  elif [[ -n "$EVM_RECEIVER_CONTRACT" ]]; then
    submit_tx "9/9: Enroll remote router (domain=$REMOTE_DOMAIN)" \
      "$BINARY" tx warp enroll-remote-router \
      "$token_id" "$REMOTE_DOMAIN" "$EVM_RECEIVER_CONTRACT" "$ENROLL_GAS"
    save_state "router_enrolled" "true"
    log_ok "9/9: Remote router enrolled for domain $REMOTE_DOMAIN"
  else
    log_warn "9/9: Skipped — no EVM address yet. Run 03-evm-deploy.sh to deploy + enroll."
  fi
}

# ─── Verify ──────────────────────────────────────────────────────────────────

verify_setup() {
  log_step "Verifying Hyperlane Setup"

  local qf=(--output json --node "$NODE" --home "$BITSONG_HOME")

  log "Mailboxes:";      "$BINARY" query hyperlane mailboxes "${qf[@]}" | jq '.' 2>/dev/null || true; echo
  log "ISMs:";           "$BINARY" query hyperlane ism isms "${qf[@]}" | jq '.' 2>/dev/null || true; echo
  log "IGPs:";           "$BINARY" query hyperlane hooks igps "${qf[@]}" | jq '.' 2>/dev/null || true; echo
  log "MerkleTree hooks:"; "$BINARY" query hyperlane hooks merkle-tree-hooks "${qf[@]}" | jq '.' 2>/dev/null || true; echo
  log "Warp tokens:";    "$BINARY" query warp tokens "${qf[@]}" | jq '.' 2>/dev/null || true; echo

  local token_id; token_id=$(load_state "token_id")
  if [[ -n "$token_id" ]]; then
    log "Remote routers for $token_id:"
    "$BINARY" query warp remote-routers "$token_id" "${qf[@]}" | jq '.' 2>/dev/null || true; echo
  fi

  # Summary
  log_step "Summary"
  echo -e "${BOLD}Chain${NC}:     $CHAIN_ID  (RPC: $NODE)"
  echo -e "${BOLD}Domain${NC}:    $DOMAIN_ID  (Remote: $REMOTE_DOMAIN)"
  echo -e "${BOLD}Mailbox${NC}:   $(load_state mailbox_id)"
  echo -e "${BOLD}ISM${NC}:       $(load_state routing_ism_id)"
  echo -e "${BOLD}Hooks${NC}:     merkle=$(load_state merkle_hook_id)  igp=$(load_state igp_id)"
  echo -e "${BOLD}Token${NC}:     $(load_state token_id)"
  echo -e "${BOLD}State${NC}:     $STATE_FILE"
}

# ─── Update Enrollment ───────────────────────────────────────────────────────

update_enrollment() {
  local evm_addr="$1"
  log_step "Updating Remote Router Enrollment"

  [[ "$evm_addr" =~ ^0x[0-9a-fA-F]{40}$ ]] || { log_err "Invalid EVM address: $evm_addr"; return 1; }

  local evm_bytes32; evm_bytes32=$(evm_to_bytes32 "$evm_addr")
  local token_id; token_id=$(load_state "token_id")
  [[ -n "$token_id" ]] || { log_err "No token_id in state. Run full setup first."; return 1; }

  log "EVM address:   $evm_addr"
  log "As bytes32:    $evm_bytes32"
  log "Token:         $token_id"

  submit_tx "Update remote router (domain=$REMOTE_DOMAIN)" \
    "$BINARY" tx warp enroll-remote-router \
    "$token_id" "$REMOTE_DOMAIN" "$evm_bytes32" "$ENROLL_GAS"

  save_state "evm_hyp_erc20" "$evm_addr"
  save_state "evm_hyp_erc20_bytes32" "$evm_bytes32"
  save_state "evm_enrollment_updated" "true"
  log_ok "Enrollment updated"

  "$BINARY" query warp remote-routers "$token_id" \
    --output json --node "$NODE" --home "$BITSONG_HOME" | jq '.' 2>/dev/null || true
}

# ─── Main ────────────────────────────────────────────────────────────────────

UPDATE_EVM_ADDR=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --update-enrollment) UPDATE_EVM_ADDR="$2"; shift 2 ;;
    -h|--help)
      echo "Usage: $0 [--update-enrollment <evm_addr>]"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

require_binary
require_jq
require_chain_running
log "Using binary: $BINARY ($("$BINARY" version 2>&1))"

banner "Hyperlane Bridge Setup" "Domain: $DOMAIN_ID"

setup_hyperlane
verify_setup

if [[ -n "$UPDATE_EVM_ADDR" ]]; then
  update_enrollment "$UPDATE_EVM_ADDR"
fi

log_ok "Done!"
