# Single-node quickstart

Use this guide to bring up the single-node stack (Kong → Portal → Core services) with Docker Compose.

## Prerequisites

- Docker and Docker Compose v2
- Make
- OpenSSL (used by the secrets script)

## Steps

1. Copy the sample environment:

   ```bash
   cp .env.example .env
   ```

2. Generate required secrets (passwords, tokens, tunnel credentials placeholder):

   ```bash
   ./scripts/setup-secrets.sh
   ```

3. Start the core + data stack:

   ```bash
   make dev-up
   ```

4. (Optional) Start observability tools:

   ```bash
   make obs-up
   ```

5. Verify key services from the host:

   ```bash
   docker compose -f docker-compose.core.yml -f docker-compose.data.yml ps
   curl -I http://localhost:8000/
   curl -I http://localhost:8000/api/health
   ```

6. Stop stacks when finished:

   ```bash
   make dev-down
   make obs-down
   ```

> Kong proxies `/` to the Next.js portal and `/api` to the Go core service. Postgres and PgBouncer stay on the shared `sunmart_net` network and are ready for future multi-node work.
