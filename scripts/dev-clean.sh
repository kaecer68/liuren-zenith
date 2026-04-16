#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PORTS_FILE="$REPO_ROOT/.env.ports"

if [[ -f "$PORTS_FILE" ]]; then
  # shellcheck disable=SC1090
  source "$PORTS_FILE"
fi

LIUREN_GRPC_PORT="${LIUREN_GRPC_PORT:-${GRPC_PORT:-}}"
LIUREN_REST_PORT="${LIUREN_REST_PORT:-${REST_PORT:-}}"

if [[ -z "$LIUREN_GRPC_PORT" || -z "$LIUREN_REST_PORT" ]]; then
  echo "[dev-clean] 缺少 runtime port 設定，請先執行 scripts/sync-contracts.sh" >&2
  exit 1
fi

ports=("$LIUREN_GRPC_PORT" "$LIUREN_REST_PORT")
for port in "${ports[@]}"; do
  pids="$(lsof -tiTCP:"$port" -sTCP:LISTEN || true)"
  if [[ -n "$pids" ]]; then
    echo "[dev-clean] 清理 port $port: $pids"
    kill $pids 2>/dev/null || true
  fi
done
