#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SECRETS_DIR="${ROOT_DIR}/secrets"
mkdir -p "${SECRETS_DIR}"

rand_b64() {
  openssl rand -base64 "$1" | tr -d '\n'
}

create_secret() {
  local file="$1"; shift
  local label="$1"; shift
  local length="${1:-32}"

  if [[ -s "${file}" ]]; then
    echo "[skip] ${label} already exists at ${file}"
    return
  fi

  rand_b64 "${length}" >"${file}"
  echo "[ok] generated ${label} at ${file}"
}

create_secret "${SECRETS_DIR}/postgres_password" "PostgreSQL password" 24

# PgBouncer uses the same backend credentials to authenticate
if [[ ! -s "${SECRETS_DIR}/pgbouncer_auth" ]]; then
  printf "sunmart:%s" "$(cat "${SECRETS_DIR}/postgres_password")" >"${SECRETS_DIR}/pgbouncer_auth"
  echo "[ok] generated PgBouncer auth record at ${SECRETS_DIR}/pgbouncer_auth"
else
  echo "[skip] PgBouncer auth already exists at ${SECRETS_DIR}/pgbouncer_auth"
fi

if [[ ! -s "${SECRETS_DIR}/portainer_admin_password" ]]; then
  ADMIN_PASS="$(rand_b64 12)"
  HASHED_PASS="$(openssl passwd -bcrypt "${ADMIN_PASS}")"
  echo "${HASHED_PASS}" >"${SECRETS_DIR}/portainer_admin_password"
  echo "[ok] generated Portainer admin password (store this!): ${ADMIN_PASS}"
else
  echo "[skip] Portainer admin password already exists"
fi

create_secret "${SECRETS_DIR}/kong_admin_password" "Kong admin token" 18
create_secret "${SECRETS_DIR}/jwt_secret" "JWT signing secret" 32

echo "Secrets ready under ${SECRETS_DIR}"
