# Single-node quickstart

This guide explains how to bring up the single-node stack with PostgreSQL, PgBouncer, NATS JetStream, Kong, the Next.js portal, and observability tools.

## Prerequisites
- Docker and Docker Compose v2
- Node.js 20+ and Go 1.22+ (for local builds/tests)

## Steps
1. Copy the example environment and adjust values as needed:
   ```bash
   cp .env.example .env
   ```

2. Generate local secrets (passwords/tokens are written under `secrets/`):
   ```bash
   ./scripts/setup-secrets.sh
   ```

3. Start the core + data stack:
   ```bash
   make dev-up
   ```

4. (Optional) Start observability services in a separate terminal:
   ```bash
   make obs-up
   ```

5. Verify the stack:
   - Kong proxy: `curl -I http://localhost:8000/`
   - Frontend health: `curl http://localhost:8000/api/health`
   - NATS monitor: `curl http://localhost:8222/healthz`
   - PostgreSQL readiness: `docker compose -f docker-compose.data.yml ps`

6. When finished, tear down the services:
   ```bash
   make dev-down
   make obs-down
   ```
