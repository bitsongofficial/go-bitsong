#!/usr/bin/env bash
# =============================================================================
# 04-agents.sh — Replace NoopISMs with real MultisigISMs, start validators
#                and relayer via Docker
#
# Usage:
#   bash 04-agents.sh                  # Full setup
#   bash 04-agents.sh --validator-only # Start validators only
#   bash 04-agents.sh --relayer-only   # Start relayer only
#   bash 04-agents.sh --clean          # Stop containers + wipe Phase 5 state
#   bash 04-agents.sh --stop           # Stop all Phase 5 containers
#
# Required environment:
#   VALIDATOR_KEY       — EVM hex private key for validator signing
#   COSMOS_SIGNER_KEY   — hex key for Cosmos announcement txs
#   EVM_RELAYER_KEY     — hex key for Base Sepolia relayer (HYP_KEY alias)
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

HYP_KEY="${HYP_KEY:-${EVM_RELAYER_KEY:-}}"
EVM_RELAYER_KEY="${EVM_RELAYER_KEY:-${HYP_KEY:-}}"

# ─── Step 0: Prerequisites ──────────────────────────────────────────────────

check_prerequisites() {
  log_step "Step 0: Prerequisites"
  local ok=true

  for cmd in docker cast jq; do
    command -v "$cmd" >/dev/null 2>&1 && log_ok "$cmd: $(command -v "$cmd")" \
      || { log_err "$cmd not found"; ok=false; }
  done

  require_binary
  log_ok "bitsongd: $BINARY ($("$BINARY" version 2>&1 || echo 'unknown'))"

  [[ -n "${VALIDATOR_KEY:-}" ]]     || { log_err "VALIDATOR_KEY not set"; ok=false; }
  [[ -n "${COSMOS_SIGNER_KEY:-}" ]] || { log_err "COSMOS_SIGNER_KEY not set"; ok=false; }
  [[ -n "${EVM_RELAYER_KEY:-}" ]]   || { log_err "EVM_RELAYER_KEY (or HYP_KEY) not set"; ok=false; }

  [[ -f "$STATE_FILE" ]] || { log_err "State file not found: $STATE_FILE"; ok=false; }
  if [[ -f "$STATE_FILE" ]]; then
    require_state "mailbox_id" "mailbox_id" > /dev/null
    require_state "token_id" "token_id" > /dev/null
    log_ok "Phase 3/4 state verified"
  fi

  [[ "$ok" == "true" ]] || { log_err "Prerequisites failed"; exit 1; }
  log_ok "All prerequisites satisfied"
}

# ─── Step 1: Derive Validator Address ────────────────────────────────────────

derive_validator_addr() {
  log_step "Step 1: Derive Validator Address"
  VALIDATOR_ADDR=$(load_state "validator_addr")
  if [[ -n "$VALIDATOR_ADDR" ]]; then
    log_ok "Validator (cached): $VALIDATOR_ADDR"
  else
    VALIDATOR_ADDR=$(cast wallet address --private-key "$VALIDATOR_KEY" 2>&1)
    [[ "$VALIDATOR_ADDR" =~ ^0x[0-9a-fA-F]{40}$ ]] || { log_err "Failed to derive validator address"; return 1; }
    save_state "validator_addr" "$VALIDATOR_ADDR"
    log_ok "Validator: $VALIDATOR_ADDR"
  fi
}

# ─── Step 1b: Fund Cosmos Signer ────────────────────────────────────────────

fund_cosmos_signer() {
  log_step "Step 1b: Fund Cosmos Signer"

  local cosmos_signer_addr
  cosmos_signer_addr=$(python3 -c "
from ecdsa import SigningKey, SECP256k1
import hashlib
def to_bech32(hex_key, prefix='bitsong'):
    pk = bytes.fromhex(hex_key.replace('0x',''))
    sk = SigningKey.from_string(pk, curve=SECP256k1)
    vk = sk.get_verifying_key()
    x, y = vk.to_string()[:32], vk.to_string()[32:]
    compressed = (b'\x02' if y[-1] % 2 == 0 else b'\x03') + x
    ripe = hashlib.new('ripemd160', hashlib.sha256(compressed).digest()).digest()
    CHARSET = 'qpzry9x8gf2tvdw0s3jn54khce6mua7l'
    def polymod(values):
        GEN = [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3]
        chk = 1
        for v in values:
            b = (chk >> 25); chk = (chk & 0x1ffffff) << 5 ^ v
            for i in range(5): chk ^= GEN[i] if ((b >> i) & 1) else 0
        return chk
    hrp_exp = [ord(x) >> 5 for x in prefix] + [0] + [ord(x) & 31 for x in prefix]
    def conv(data):
        acc, bits, ret = 0, 0, []
        for value in data:
            acc = (acc << 8) | value; bits += 8
            while bits >= 5: bits -= 5; ret.append((acc >> bits) & 31)
        if bits: ret.append((acc << (5 - bits)) & 31)
        return ret
    data = conv(list(ripe))
    ck = [(polymod(hrp_exp + data + [0]*6) ^ 1) >> 5*(5-i) & 31 for i in range(6)]
    return prefix + '1' + ''.join([CHARSET[d] for d in data + ck])
print(to_bech32('$COSMOS_SIGNER_KEY'))
" 2>/dev/null)

  if [[ -z "$cosmos_signer_addr" || ! "$cosmos_signer_addr" =~ ^bitsong1 ]]; then
    log_err "Failed to derive Cosmos signer address. Ensure: pip3 install ecdsa"
    return 1
  fi

  log "Cosmos signer: $cosmos_signer_addr"

  local balance
  balance=$("$BINARY" query bank balances "$cosmos_signer_addr" --home "$BITSONG_HOME" --output json 2>/dev/null \
    | jq -r '.balances[] | select(.denom=="ubtsg") | .amount // "0"') || balance="0"
  if [[ -n "$balance" && "$balance" != "0" && "$balance" != "null" ]]; then
    log_ok "Already funded: ${balance}ubtsg"
    return 0
  fi

  log "Funding with 10000000ubtsg..."
  submit_tx "Fund cosmos signer" "$BINARY" tx bank send val "$cosmos_signer_addr" 10000000ubtsg
  log_ok "Cosmos signer funded"
}

# ─── Step 2: Update ISM — Cosmos Side ───────────────────────────────────────

update_ism_cosmos_side() {
  log_step "Step 2: Update ISM — Cosmos Side (EVM->Cosmos)"

  [[ "$(load_state "ism_updated_cosmos_side")" == "true" ]] && { log_ok "Already done"; return 0; }

  local routing_ism_id; routing_ism_id=$(require_state "routing_ism_id" "routing_ism_id")

  # Create MultisigISM with real validator (1-of-1)
  local bitsong_multisig_ism_id
  bitsong_multisig_ism_id=$(load_state "bitsong_multisig_ism_id")
  if [[ -z "$bitsong_multisig_ism_id" ]]; then
    submit_tx "Create MultisigISM (validator=$VALIDATOR_ADDR, threshold=1)" \
      "$BINARY" tx hyperlane ism create-message-id-multisig "$VALIDATOR_ADDR" "1"
    bitsong_multisig_ism_id=$(extract_id)
    [[ -n "$bitsong_multisig_ism_id" ]] || { log_err "Failed to extract ISM ID"; return 1; }
    save_state "bitsong_multisig_ism_id" "$bitsong_multisig_ism_id"
    log_ok "New MultisigISM: $bitsong_multisig_ism_id"
  else
    log_ok "MultisigISM exists: $bitsong_multisig_ism_id"
  fi

  # Remove-then-add pattern (avoids Go range-copy bug in SetDomain)
  submit_tx "Remove domain $REMOTE_DOMAIN from RoutingISM" \
    "$BINARY" tx hyperlane ism remove-routing-ism-domain "$routing_ism_id" "$REMOTE_DOMAIN"
  submit_tx "Set domain $REMOTE_DOMAIN to new MultisigISM" \
    "$BINARY" tx hyperlane ism set-routing-ism-domain "$routing_ism_id" "$REMOTE_DOMAIN" "$bitsong_multisig_ism_id"

  save_state "ism_updated_cosmos_side" "true"
  log_ok "Cosmos ISM: routing[$REMOTE_DOMAIN] -> $bitsong_multisig_ism_id"
}

# ─── Step 3: Update ISM — EVM Side ──────────────────────────────────────────

update_ism_evm_side() {
  log_step "Step 3: Update ISM — EVM Side (Cosmos->EVM)"

  [[ "$(load_state "ism_updated_evm_side")" == "true" ]] && { log_ok "Already done"; return 0; }

  local evm_hyp_erc20; evm_hyp_erc20=$(require_state "evm_hyp_erc20" "evm_hyp_erc20")

  # Deploy ISM via factory (CREATE2 — deterministic)
  local basesepolia_multisig_ism
  basesepolia_multisig_ism=$(load_state "basesepolia_multisig_ism")
  if [[ -z "$basesepolia_multisig_ism" ]]; then
    log "Deploying staticMessageIdMultisigISM on Base Sepolia..."

    basesepolia_multisig_ism=$(cast call "$BASESEPOLIA_ISM_FACTORY" \
      "deploy(address[],uint8)(address)" "[$VALIDATOR_ADDR]" 1 \
      --rpc-url "$EVM_RPC" 2>/dev/null) || true
    [[ "$basesepolia_multisig_ism" =~ ^0x[0-9a-fA-F]{40}$ ]] || { log_err "Failed to predict ISM address"; return 1; }
    log "  Predicted: $basesepolia_multisig_ism"

    local tx_out; tx_out=$(cast send "$BASESEPOLIA_ISM_FACTORY" \
      "deploy(address[],uint8)" "[$VALIDATOR_ADDR]" 1 \
      --private-key "$HYP_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true
    local tx_hash; tx_hash=$(echo "$tx_out" | jq -r '.transactionHash // empty' 2>/dev/null) || true
    [[ -n "$tx_hash" ]] || { log_err "ISM deploy failed"; echo "$tx_out"; return 1; }

    save_state "basesepolia_multisig_ism" "$basesepolia_multisig_ism"
    log_ok "ISM deployed: $basesepolia_multisig_ism (tx: $tx_hash)"
  else
    log_ok "ISM exists: $basesepolia_multisig_ism"
  fi

  # Wait for ISM tx confirmation, then update HypERC20
  sleep 5
  log "Setting HypERC20 ISM to $basesepolia_multisig_ism..."
  local update_out; update_out=$(cast send "$evm_hyp_erc20" \
    "setInterchainSecurityModule(address)" "$basesepolia_multisig_ism" \
    --private-key "$HYP_KEY" --rpc-url "$EVM_RPC" --json 2>&1) || true
  local update_status; update_status=$(echo "$update_out" | jq -r '.status // empty' 2>/dev/null) || true
  [[ "$update_status" == "0x1" || "$update_status" == "1" ]] \
    || { log_err "setInterchainSecurityModule failed (status: $update_status)"; echo "$update_out"; return 1; }

  save_state "ism_updated_evm_side" "true"
  log_ok "HypERC20 ISM updated to real MultisigISM"
}

# ─── Agent Config ────────────────────────────────────────────────────────────

write_agent_config() {
  local mailbox_id merkle_hook_id igp_id
  mailbox_id=$(load_state "mailbox_id")
  merkle_hook_id=$(load_state "merkle_hook_id")
  igp_id=$(load_state "igp_id")

  # CRITICAL: index.from for basesepolia MUST be the mailbox deployment block (~13,850,000).
  # The relayer's msg::db_loader scans by mailbox nonce starting from 0 sequentially.
  # If index.from is too recent, only high nonces are indexed and db_loader is stuck at 0.
  local basesep_index_from=13850000

  cat > "$BITSONG_HOME/agent-config.json" << EOF
{
  "chains": {
    "bitsong": {
      "name": "bitsong",
      "chainId": "${CHAIN_ID}",
      "domainId": $DOMAIN_ID,
      "protocol": "cosmosNative",
      "bech32Prefix": "bitsong",
      "slip44": 639,
      "contractAddressBytes": 32,
      "canonicalAsset": "ubtsg",
      "rpcUrls": [{"http": "http://localhost:26657"}],
      "grpcUrls": [{"http": "http://localhost:9090"}],
      "nativeToken": { "name": "BitSong", "symbol": "BTSG", "decimals": 6, "denom": "ubtsg" },
      "gasPrice": { "amount": "0.025", "denom": "ubtsg" },
      "gasMultiplier": "2.0",
      "blocks": { "confirmations": 1, "estimateBlockTime": 6, "reorgPeriod": 1 },
      "index": { "from": 1, "chunk": 50 },
      "mailbox": "${mailbox_id}",
      "validatorAnnounce": "${mailbox_id}",
      "merkleTreeHook": "${merkle_hook_id}",
      "interchainGasPaymaster": "${igp_id}"
    },
    "basesepolia": {
      "name": "basesepolia",
      "chainId": 84532,
      "domainId": 84532,
      "protocol": "ethereum",
      "rpcUrls": [{"http": "${EVM_RPC}"}],
      "nativeToken": { "name": "Ether", "symbol": "ETH", "decimals": 18 },
      "blocks": { "confirmations": 1, "estimateBlockTime": 2, "reorgPeriod": 1 },
      "index": { "from": ${basesep_index_from}, "chunk": 9999 },
      "mailbox": "0x6966b0E55883d49BFB24539356a2f8A673E02039",
      "validatorAnnounce": "${VALIDATOR_ANNOUNCE_CONTRACT}",
      "merkleTreeHook": "0x86fb9F1c124fB20ff130C41a79a432F770f67AFD",
      "interchainGasPaymaster": "0x28B02B97a850872C4D33C3E024fab6499ad96564",
      "interchainSecurityModule": "0xBB276c7419155980558BFf56E22AfF83023a2dB2",
      "proxyAdmin": "0x44b764045BfDC68517e10e783E69B376cef196B2"
    }
  }
}
EOF
  log_ok "Agent config: $BITSONG_HOME/agent-config.json (basesepolia index.from=$basesep_index_from)"
}

# ─── Step 4: Start Validator — BitSong ───────────────────────────────────────

start_validator_bitsong() {
  log_step "Step 4: Start Validator — BitSong"

  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^hyperlane-validator-bitsong$"; then
    log_ok "Already running"; return 0
  fi
  docker rm -f hyperlane-validator-bitsong 2>/dev/null || true

  mkdir -p "$BITSONG_HOME/validator-bitsong-db" "$BITSONG_HOME/checkpoints-bitsong"
  write_agent_config

  log "Starting validator-bitsong (domain=$DOMAIN_ID)..."
  docker run -d --name hyperlane-validator-bitsong --network host \
    -e CONFIG_FILES=/config/agent-config.json \
    -v "$BITSONG_HOME/agent-config.json:/config/agent-config.json:ro" \
    -v "$BITSONG_HOME/validator-bitsong-db:/hyperlane_db" \
    -v "$BITSONG_HOME/checkpoints-bitsong:/checkpoints-bitsong" \
    "$DOCKER_IMAGE" ./validator \
    --db /hyperlane_db --originChainName bitsong \
    --reorgPeriod 1 --interval 5 \
    --validator.type hexKey --validator.key "$VALIDATOR_KEY" \
    --chains.bitsong.signer.type cosmosKey \
    --chains.bitsong.signer.key "$COSMOS_SIGNER_KEY" \
    --chains.bitsong.signer.prefix bitsong \
    --checkpointSyncer.type localStorage \
    --checkpointSyncer.path /checkpoints-bitsong

  log_ok "validator-bitsong started"
}

# ─── Step 5: Start Validator — Base Sepolia ──────────────────────────────────

start_validator_basesepolia() {
  log_step "Step 5: Start Validator — Base Sepolia"

  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^hyperlane-validator-basesepolia$"; then
    log_ok "Already running"; return 0
  fi
  docker rm -f hyperlane-validator-basesepolia 2>/dev/null || true

  mkdir -p "$BITSONG_HOME/validator-basesepolia-db" "$BITSONG_HOME/checkpoints-basesepolia"

  log "Starting validator-basesepolia (domain=$REMOTE_DOMAIN)..."
  docker run -d --name hyperlane-validator-basesepolia --network host \
    -e CONFIG_FILES=/config/agent-config.json \
    -v "$BITSONG_HOME/agent-config.json:/config/agent-config.json:ro" \
    -v "$BITSONG_HOME/validator-basesepolia-db:/hyperlane_db" \
    -v "$BITSONG_HOME/checkpoints-basesepolia:/checkpoints-basesepolia" \
    "$DOCKER_IMAGE" ./validator \
    --db /hyperlane_db --originChainName basesepolia \
    --reorgPeriod 1 --interval 5 \
    --validator.type hexKey --validator.key "$VALIDATOR_KEY" \
    --chains.basesepolia.signer.type hexKey \
    --chains.basesepolia.signer.key "$EVM_RELAYER_KEY" \
    --checkpointSyncer.type localStorage \
    --checkpointSyncer.path /checkpoints-basesepolia

  log_ok "validator-basesepolia started"
}

# ─── Step 6: Wait for Announcements ─────────────────────────────────────────

wait_for_announcements() {
  log_step "Step 6: Wait for Validator Announcements"

  local bitsong_ok basesepolia_ok
  bitsong_ok=$(load_state "validator_bitsong_announced")
  basesepolia_ok=$(load_state "validator_basesepolia_announced")
  [[ "$bitsong_ok" == "true" && "$basesepolia_ok" == "true" ]] && { log_ok "Both already announced"; return 0; }

  log "Waiting up to 120s..."
  local validator_addr_lower; validator_addr_lower=$(echo "$VALIDATOR_ADDR" | tr '[:upper:]' '[:lower:]')

  for i in $(seq 1 24); do
    sleep 5

    if [[ "$bitsong_ok" != "true" ]]; then
      local mailbox_id locations
      mailbox_id=$(load_state "mailbox_id")
      locations=$("$BINARY" query hyperlane ism announced-storage-locations \
        "$mailbox_id" "$validator_addr_lower" \
        --output json --node "$NODE" --home "$BITSONG_HOME" 2>/dev/null | \
        jq -r '.storage_locations // [] | length' 2>/dev/null) || locations="0"
      if [[ "$locations" -gt 0 ]]; then
        save_state "validator_bitsong_announced" "true"; bitsong_ok="true"
        log_ok "validator-bitsong announced"
      fi
    fi

    if [[ "$basesepolia_ok" != "true" ]]; then
      local evm_announced
      evm_announced=$(cast call "$VALIDATOR_ANNOUNCE_CONTRACT" \
        "getAnnouncedValidators()(address[])" --rpc-url "$EVM_RPC" 2>/dev/null | tr ',' '\n' | \
        grep -i "$(echo "${VALIDATOR_ADDR#0x}" | tr '[:upper:]' '[:lower:]')" || true)
      if [[ -n "$evm_announced" ]]; then
        save_state "validator_basesepolia_announced" "true"; basesepolia_ok="true"
        log_ok "validator-basesepolia announced"
      fi
    fi

    [[ "$bitsong_ok" == "true" && "$basesepolia_ok" == "true" ]] && return 0
    log "  [${i}/24] bitsong=${bitsong_ok:-pending}, basesepolia=${basesepolia_ok:-pending}"
  done

  log_warn "Announcement timeout — validators may still be starting. Check: docker logs hyperlane-validator-bitsong"
}

# ─── Step 7: Start Relayer ──────────────────────────────────────────────────

start_relayer() {
  log_step "Step 7: Start Relayer"

  if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "^hyperlane-relayer$"; then
    log_ok "Already running"; return 0
  fi
  docker rm -f hyperlane-relayer 2>/dev/null || true

  mkdir -p "$BITSONG_HOME/relayer-db"

  # Whitelist: only relay our warp messages (avoids noise from other projects on shared testnet)
  local token_id evm_hyp_erc20 evm_padded
  token_id=$(load_state "token_id")
  evm_hyp_erc20=$(load_state "evm_hyp_erc20")
  evm_padded=$(evm_to_bytes32 "$evm_hyp_erc20")

  local whitelist="[{\"senderAddress\":\"${token_id}\",\"destinationDomain\":\"${REMOTE_DOMAIN}\"},{\"senderAddress\":\"${evm_padded}\",\"destinationDomain\":\"${DOMAIN_ID}\"}]"

  log "Starting relayer (bitsong <-> basesepolia)..."
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

  log_ok "Relayer started"
}

# ─── Main ────────────────────────────────────────────────────────────────────

CLEAN=false
STOP=false
VALIDATOR_ONLY=false
RELAYER_ONLY=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --clean)          CLEAN=true; shift ;;
    --stop)           STOP=true; shift ;;
    --validator-only) VALIDATOR_ONLY=true; shift ;;
    --relayer-only)   RELAYER_ONLY=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--clean|--stop|--validator-only|--relayer-only]"
      echo ""
      echo "Required env: VALIDATOR_KEY, COSMOS_SIGNER_KEY, EVM_RELAYER_KEY"
      exit 0 ;;
    *) log_err "Unknown flag: $1"; exit 1 ;;
  esac
done

if [[ "$STOP" == "true" ]]; then
  log "Stopping Phase 5 containers..."
  docker stop hyperlane-validator-bitsong hyperlane-validator-basesepolia hyperlane-relayer 2>/dev/null || true
  log_ok "Stopped"; exit 0
fi

if [[ "$CLEAN" == "true" ]]; then
  log "Cleaning Phase 5 state..."
  docker rm -f hyperlane-validator-bitsong hyperlane-validator-basesepolia hyperlane-relayer 2>/dev/null || true
  for key in validator_addr bitsong_multisig_ism_id basesepolia_multisig_ism \
             validator_bitsong_announced validator_basesepolia_announced \
             ism_updated_cosmos_side ism_updated_evm_side; do
    [[ -f "$STATE_FILE" ]] && jq --arg k "$key" 'del(.[$k])' "$STATE_FILE" > "${STATE_FILE}.tmp" \
      && mv "${STATE_FILE}.tmp" "$STATE_FILE"
  done
  log_ok "Phase 5 state cleaned"
fi

banner "Phase 5: Agents" "bitsong <-> basesepolia"

check_prerequisites
derive_validator_addr
fund_cosmos_signer

if [[ "$RELAYER_ONLY" != "true" ]]; then
  update_ism_cosmos_side
  update_ism_evm_side
  start_validator_bitsong
  start_validator_basesepolia
  wait_for_announcements
fi

if [[ "$VALIDATOR_ONLY" != "true" ]]; then
  start_relayer
fi

log_ok "Agents running! Next: bash 05-test.sh"
