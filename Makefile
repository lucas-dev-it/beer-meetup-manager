ROOT_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
DOCKER_DIR := $(ROOT_DIR)
DOCKER_FILE := $(DOCKER_DIR)/docker-compose.yml

.DEFAULT_TARGET: docker

.PHONY: run-all
run-all:
	@echo "Starting whole environment"
	@docker-compose -f $(DOCKER_FILE) --project-directory $(DOCKER_DIR) up -d --force-recreate --build --remove-orphans

.PHONY: run-app
run-app:
	@echo "Starting service instance"
	@docker-compose -f $(DOCKER_FILE) --project-directory $(DOCKER_DIR) up -d --build --remove-orphans meetup

.PHONY: docker
docker: prepare-environment
	@echo "Creating services instances"
	@docker-compose -f $(DOCKER_FILE) --project-directory $(DOCKER_DIR) up -d --build --remove-orphans redis postgres

.PHONY: prepare-environment
prepare-environment:
	@DB_USER=$(or $(DB_USER), $(shell read -p "Postgres User: " user; echo "DB_USER="$$user)); \
	DB_PASS=$(or $(DB_PASS), $(shell read -p "Postgres Pass: " pass; echo "DB_PASS="$$pass)); \
	DB_NAME=$(or $(DB_NAME), $(shell read -p "Postgres DB Name: " name; echo "DB_NAME="$$name)); \
	WEATHERBIT_API_KEY=$(or $(WEATHERBIT_API_KEY), $(shell read -p "WeatherBit API Key: " weather_bit_api_key; echo "WEATHERBIT_API_KEY="$$weather_bit_api_key)); \
	TOKEN_SIGNING_KEY=$(or $(TOKEN_SIGNING_KEY), $(shell read -p "JWT Token Signing Key: " TOKEN_SIGNING_KEY; echo "TOKEN_SIGNING_KEY="$$TOKEN_SIGNING_KEY)); \
	printf "$$DB_USER\n$$DB_PASS\n$$DB_NAME\n$$WEATHERBIT_API_KEY\n$$TOKEN_SIGNING_KEY\n" > .env

.PHONY: docker-down
docker-down:
	docker-compose -f $(DOCKER_FILE) down -v

.PHONY: redis-get-keys
redis-get-keys:
	@KEY=$(or $(KEY), $(shell read -p "Key name: " key; echo $$key)); \
	docker-compose -f $(DOCKER_FILE) exec redis \
		redis-cli KEYS "$$KEY*"

.PHONY: redis-get
redis-get:
	@KEY=$(or $(KEY), $(shell read -p "Key: " key; echo $$key)); \
	docker-compose -f $(DOCKER_FILE) exec redis \
		redis-cli GET $$KEY