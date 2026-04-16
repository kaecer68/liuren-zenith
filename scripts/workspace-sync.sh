#!/usr/bin/env bash
#
# Thin wrapper for central workspace sync
# Calls destiny-cloud's central sync engine
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WORKSPACE_ROOT="$(cd "$REPO_ROOT/.." && pwd)"
DESTINY_CLOUD="$WORKSPACE_ROOT/destiny-cloud"

# Check if destiny-cloud workspace-sync exists
if [[ ! -f "$DESTINY_CLOUD/bin/workspace-sync" ]]; then
  echo "[error] Central workspace-sync not found. Building..." >&2
  (cd "$DESTINY_CLOUD" && go build -o bin/workspace-sync ./cmd/workspace-sync)
fi

# Forward to central sync engine
exec "$DESTINY_CLOUD/bin/workspace-sync" \
  -registry "$DESTINY_CLOUD/configs/workspace/services.yaml" \
  -contracts "$WORKSPACE_ROOT/destiny-contracts" \
  -workspace "$WORKSPACE_ROOT" \
  -service liuren-zenith \
  "$@"
