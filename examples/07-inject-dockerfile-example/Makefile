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
	@echo " [$(EXAMPLE)]$(GREEN) Create SSH keys$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) run --rm --workdir /ssh ssh-keygen /usr/bin/ssh-keygen -t rsa -q -N "password" -f id_rsa -C "apenella@stevedore.test" || { echo " Error creating the SSH keys $(EXAMPLE)"; exit 1; }
	@echo

	@echo
	@echo "$(GREEN)Starting the stack to run $(EXAMPLE)$(NC)"
	@echo
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) up --detach --build || { echo " Error starting the stack to run $(EXAMPLE)"; exit 1; }

	@echo " [$(EXAMPLE)]$(GREEN) Create known_hosts file$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec --workdir /ssh ssh-keygen bash -c '/usr/bin/ssh-keyscan -H gitserver.stevedore.test > known_hosts' || { echo " Error creating the known_host keys $(EXAMPLE)"; exit 1; }
	@echo

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

example: ## Executes the example. It expects that the required stack is already running
	@echo
	@echo "$(GREEN) Run example $(EXAMPLE)$(NC)"
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Sync$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore /usr/local/bin/wait-for-dockerd.sh || { echo " Error synchronizing $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Get credentials$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore get credentials || { echo " Error getting credentials $(EXAMPLE)"; exit 1; }
	@echo

	@echo " [$(EXAMPLE)]$(GREEN) Get images$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore get images --tree || { echo " Error getting images $(EXAMPLE)"; exit 1; }
	@echo

# @echo " [$(EXAMPLE)]$(GREEN) Build the base image and its descendants, and push the images after build$(NC)"
# @$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore docker info || { echo " Error runnig example $(EXAMPLE)"; exit 1; }
# @echo

	@echo " [$(EXAMPLE)]$(GREEN) Build the base image and its descendants, and push the images after build$(NC)"
	@$(DOCKER_COMPOSE_VERSION) -p $(EXAMPLE) exec stevedore stevedore build base --build-on-cascade --push-after-build || { echo " Error runnig example $(EXAMPLE)"; exit 1; }
	@echo

a: b c d
b:
	@echo b
c:
	@echo c
d:
	@echo d