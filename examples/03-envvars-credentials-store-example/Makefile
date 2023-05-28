MAKEFLAGS += --keep-going

DOCKER_COMPOSE_VERSION := $(shell docker compose version > /dev/null 2>&1 && echo "docker compose" || (which docker-compose > /dev/null 2>&1 && echo "docker-compose" || (echo "docker compose not found. Aborting." >&2; exit 1)))

EXAMPLE := $(shell basename ${PWD})
GREEN := "\\033[32m"
NC := "\\033[0m"

.PHONY: start run clean example

help: ## Lists allowed targets
	@echo "$(GREEN) Example $(EXAMPLE):$(NC)"
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-25s\033[0m %s\n", $$1, $$2}'
	@echo 

run: start example clean ## Runs start, example, and clean targets together

start: ## Starts the stack required to run the example
	@echo
	@echo "$(GREEN)Starting the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) up --detach --build || { echo " Error starting the stack to run $(EXAMPLE)"; exit 1; }


clean: ## Stops the stack required to run the example
	@echo
	@echo "$(GREEN)Stopping the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) down --volumes --remove-orphans --timeout 3 || true

status: ## Displays the status of the stack
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) ps || true

follow-logs: ## Shows the stack logs in follow mode
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) logs -f

attach: ## Attaches to the Stevedore container
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore sh

example: ## Executes the example (requires the stack to be started)
	@echo
	@echo "$(GREEN) Run example $(EXAMPLE)$(NC)"
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Sync$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore /usr/local/bin/wait-for-dockerd.sh || { echo " Error synchronizing $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Get credentials$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore get credentials --show-secrets || { echo " Error getting credentials $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Build my-app and push images after build$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore build my-app --image-version 3.2.1 --push-after-build || { echo " Error runnig example $(EXAMPLE)"; exit 1; }
	@echo " [$(EXAMPLE)]$(GREEN) Get credentials$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore get credentials --show-secrets || { echo " Error getting credentials $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Build my-app and push images after build$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore build my-app --image-version 3.2.1 --push-after-build || { echo " Error runnig example $(EXAMPLE)"; exit 1; }
	@echo
