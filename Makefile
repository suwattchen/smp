SHELL := /bin/bash

COMPOSE_CORE := docker compose -f docker-compose.core.yml -f docker-compose.data.yml
COMPOSE_OBS := docker compose -f docker-compose.obs.yml

.PHONY: dev-up dev-down obs-up obs-down logs-core

dev-up:
$(COMPOSE_CORE) up -d

dev-down:
$(COMPOSE_CORE) down

obs-up:
$(COMPOSE_OBS) up -d

obs-down:
$(COMPOSE_OBS) down

logs-core:
$(COMPOSE_CORE) logs -f core-go kong
