#!/usr/bin/env bash
# =============================================================================
# stop.sh — Stop BitSong chain and/or Hyperlane Docker agents
#
# Usage:
#   bash stop.sh                # Stop everything
#   bash stop.sh --chain-only   # Stop chain only
#   bash stop.sh --agents-only  # Stop Docker agents only
#   bash stop.sh --clean        # Stop everything + wipe all data
# =============================================================================

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

CLEAN=false
CHAIN_ONLY=false
AGENTS_ONLY=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --clean)       CLEAN=true; shift ;;
    --chain-only)  CHAIN_ONLY=true; shift ;;
    --agents-only) AGENTS_ONLY=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--clean] [--chain-only] [--agents-only]"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

# ─── Stop Docker Agents ──────────────────────────────────────────────────────

if [[ "$CHAIN_ONLY" == "false" ]]; then
  log "Stopping Hyperlane Docker agents..."
  for name in hyperlane-validator-bitsong hyperlane-validator-basesepolia hyperlane-relayer; do
    if docker ps -a --format '{{.Names}}' 2>/dev/null | grep -q "^${name}$"; then
      docker stop "$name" 2>/dev/null || true
      docker rm -f "$name" 2>/dev/null || true
      log_ok "Stopped $name"
    fi
  done
fi

# ─── Stop Chain ──────────────────────────────────────────────────────────────

if [[ "$AGENTS_ONLY" == "false" ]]; then
  local_pid_file="$BITSONG_HOME/chain.pid"
  if [[ -f "$local_pid_file" ]]; then
    PID=$(cat "$local_pid_file")
    if kill -0 "$PID" 2>/dev/null; then
      log "Stopping chain (PID=$PID)..."
      kill "$PID" 2>/dev/null || true
      for _ in $(seq 1 20); do
        kill -0 "$PID" 2>/dev/null || break
        sleep 0.5
      done
      if kill -0 "$PID" 2>/dev/null; then
        log_warn "Force killing chain..."
        kill -9 "$PID" 2>/dev/null || true
      fi
      log_ok "Chain stopped"
    else
      log "Chain not running (stale PID=$PID)"
    fi
    rm -f "$local_pid_file"
  else
    log "No chain PID file found"
  fi
fi

# ─── Clean ───────────────────────────────────────────────────────────────────

if [[ "$CLEAN" == "true" ]]; then
  log "Removing all data at $BITSONG_HOME..."
  rm -rf "$BITSONG_HOME"
  log_ok "Clean complete"
fi
