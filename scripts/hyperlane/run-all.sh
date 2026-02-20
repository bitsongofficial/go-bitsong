#!/usr/bin/env bash
# =============================================================================
# run-all.sh — Run all Hyperlane phases in sequence
#
# Usage:
#   bash run-all.sh                  # Run phases 1-5
#   bash run-all.sh --from 3         # Start from phase 3
#   bash run-all.sh --skip-test      # Skip transfer tests
#   bash run-all.sh --clean          # Clean everything first
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

FROM_PHASE=1
SKIP_TEST=false
CLEAN=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --from)      FROM_PHASE="$2"; shift 2 ;;
    --skip-test) SKIP_TEST=true; shift ;;
    --clean)     CLEAN=true; shift ;;
    -h|--help)
      echo "Usage: $0 [--from N] [--skip-test] [--clean]"
      echo ""
      echo "Phases:"
      echo "  1  Init + start local chain        (01-chain.sh)"
      echo "  2  Configure Hyperlane bridge       (02-hyperlane.sh)"
      echo "  3  Deploy HypERC20 + enrollment     (03-evm-deploy.sh)"
      echo "  4  ISM upgrade + validators/relayer (04-agents.sh)"
      echo "  5  Transfer tests                   (05-test.sh)"
      echo "  6  Fantoken warp route              (06-fantoken-route.sh)"
      echo "  7  Fantoken transfer tests          (07-fantoken-test.sh)"
      echo ""
      echo "Environment (Phase 6/7):"
      echo "  FT_SYMBOL   Fantoken symbol (default: clay)"
      echo "  FT_NAME     Fantoken name   (default: Clay Token)"
      exit 0 ;;
    *) echo "Unknown flag: $1"; exit 1 ;;
  esac
done

run_phase() {
  local num="$1" script="$2"
  shift 2
  if [[ "$FROM_PHASE" -le "$num" ]]; then
    echo ""
    echo "================================================================"
    echo "  Phase $num: $script"
    echo "================================================================"
    echo ""
    bash "$SCRIPT_DIR/$script" "$@"
  fi
}

CLEAN_FLAG=""
[[ "$CLEAN" == "true" ]] && CLEAN_FLAG="--clean"

run_phase 1 "01-chain.sh" $CLEAN_FLAG
run_phase 2 "02-hyperlane.sh"
run_phase 3 "03-evm-deploy.sh" $CLEAN_FLAG
run_phase 4 "04-agents.sh" $CLEAN_FLAG

if [[ "$SKIP_TEST" != "true" ]]; then
  run_phase 5 "05-test.sh"
fi

FT_SYMBOL="${FT_SYMBOL:-clay}"
FT_NAME="${FT_NAME:-Clay Token}"

run_phase 6 "06-fantoken-route.sh" --symbol "$FT_SYMBOL" --name "$FT_NAME" $CLEAN_FLAG

if [[ "$SKIP_TEST" != "true" ]]; then
  run_phase 7 "07-fantoken-test.sh" --symbol "$FT_SYMBOL"
fi

echo ""
echo "All phases complete!"
