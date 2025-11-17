#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SECRETS_DIR="$ROOT_DIR/secrets"

mkdir -p "$SECRETS_DIR"

create_secret() {
  local file="$SECRETS_DIR/$1"
  if [[ -s "$file" ]]; then
    echo "[skip] $1 already exists"
    return
  fi
  if [[ $# -gt 1 ]]; then
    printf "%s" "$2" >"$file"
  else
    openssl rand -hex 24 >"$file"
  fi
  chmod 600 "$file"
  echo "[ok] generated $1"
}

# Database credentials
create_secret "postgres_password"
create_secret "pgbouncer_password"

# Kong admin token
create_secret "kong_admin_token"

# Portainer admin password (set a deterministic value for local dev if provided via env)
if [[ -n "${PORTAINER_ADMIN_PASSWORD:-}" ]]; then
  create_secret "portainer_admin_password" "$PORTAINER_ADMIN_PASSWORD"
else
  create_secret "portainer_admin_password"
fi

# Prometheus basic auth password
create_secret "prometheus_basic_auth_password"

# Cloudflared credentials placeholder (users should replace with real tunnel creds)
if [[ ! -s "$SECRETS_DIR/cloudflared_credentials" ]]; then
  cat >"$SECRETS_DIR/cloudflared_credentials" <<'CLOUDFLARE_PLACEHOLDER'
{
  "AccountTag": "update-me",
  "TunnelSecret": "update-me",
  "TunnelID": "update-me"
}
CLOUDFLARE_PLACEHOLDER
  chmod 600 "$SECRETS_DIR/cloudflared_credentials"
  echo "[ok] stubbed cloudflared_credentials"
else
  echo "[skip] cloudflared_credentials already exists"
fi

echo "Secrets prepared in $SECRETS_DIR"
