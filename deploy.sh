#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="${ROOT_DIR}/deploy"

usage() {
  cat <<'EOF'
Usage: ./deploy.sh <command>

Commands:
  deploy   Pull latest code, vendor deps, rebuild and start services
  stop     Stop services
  restart  Stop then deploy
EOF
}

ensure_deploy_dir() {
  if [[ ! -d "${DEPLOY_DIR}" ]]; then
    echo "deploy directory not found: ${DEPLOY_DIR}" >&2
    exit 1
  fi
}

cmd_deploy() {
  echo "==> Updating code"
  git -C "${ROOT_DIR}" pull --rebase

  echo "==> Vendoring dependencies"
  (cd "${ROOT_DIR}" && go mod vendor)

  echo "==> Rebuilding and starting services"
  (cd "${DEPLOY_DIR}" && docker compose up -d --build)
}

cmd_stop() {
  ensure_deploy_dir
  echo "==> Stopping services"
  (cd "${DEPLOY_DIR}" && docker compose down)
}

cmd_status() {
  ensure_deploy_dir
  echo "==> Checking service status"
  (cd "${DEPLOY_DIR}" && docker compose ps)
}

cmd_restart() {
  cmd_stop
  cmd_deploy
}

case "${1:-}" in
  deploy)
    cmd_deploy
    ;;
  stop)
    cmd_stop
    ;;
  status)
    cmd_status
    ;;
  restart)
    cmd_restart
    ;;
  -h|--help|"")
    usage
    exit 0
    ;;
  *)
    echo "Unknown command: ${1}" >&2
    usage
    exit 1
    ;;
esac
