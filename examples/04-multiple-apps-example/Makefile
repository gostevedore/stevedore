MAKEFLAGS += --keep-going

DOCKER_COMPOSE_VERSION := $(shell docker compose version > /dev/null 2>&1 && echo "docker compose" || (which docker-compose > /dev/null 2>&1 && echo "docker-compose" || (echo "docker compose not found. Aborting." >&2; exit 1)))

EXAMPLE := $(shell basename ${PWD})
GREEN := "\\033[32m"
NC := "\\033[0m"

.PHONY: start run clean example

help: ## list allowed targets
	@echo "$(GREEN) Example $(EXAMPLE):$(NC)"
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-25s\033[0m %s\n", $$1, $$2}'
	@echo 

color:


run: start example clean ## Run the example

start: ## Starts the stack for required to run the example
	@echo
	@echo "$(GREEN)Starting the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) up --detach --build || { echo " Error starting the stack to run $(EXAMPLE)"; exit 1; }


clean: ## Stops the stack required to run the example
	@echo
	@echo "$(GREEN)Stopping the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) down --volumes --remove-orphans --timeout 3 || true

status: ## Show the status of the stack
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) ps || true

follow-logs: ## Show the stack logs
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) logs -f

attach: ## Attach the the stevedore container
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore sh

example: ## Executes the example. It expects that the required stack is already running
	@echo " [$(EXAMPLE)]$(GREEN) Sync$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore /usr/local/bin/wait-for-dockerd.sh || { echo " Error synchronizing $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Get images$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore get images --tree || { echo " Error getting images $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Build app1 and push images after build$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore build app1 --push-after-build || { echo " Error building app1 $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Build app2 and push images after build$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore build app2 --push-after-build || { echo " Error building app2 $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Build app3 and push images after build$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore build app3 --push-after-build || { echo " Error building app3 $(EXAMPLE)"; exit 1; }
	@echo