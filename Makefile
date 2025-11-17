COMPOSE_CORE=-f docker-compose.core.yml -f docker-compose.data.yml
COMPOSE_OBS=-f docker-compose.obs.yml

.PHONY: dev-up dev-down obs-up obs-down logs-core edge-up edge-down

dev-up:
	docker compose $(COMPOSE_CORE) up -d

dev-down:
	docker compose $(COMPOSE_CORE) down

obs-up:
	docker compose $(COMPOSE_OBS) up -d

obs-down:
	docker compose $(COMPOSE_OBS) down

logs-core:
	docker compose $(COMPOSE_CORE) logs -f core-go kong nats frontend postgres pgbouncer

edge-up:
	docker compose -f docker-compose.edge.yml up -d

edge-down:
	docker compose -f docker-compose.edge.yml down
