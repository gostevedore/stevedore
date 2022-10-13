help: ## list allowed targets
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[32m%-20s\033[0m %s\n", $$1, $$2}'
	@echo 

start: cleanup generate-keys generate-certs ## start docker registry
	docker compose -p stevedore-functional-test-manual up -d --build registry docker-hub gitserver stevedore

prepare-images:
	docker compose -p stevedore-functional-test-manual run stevedore /prepare-images.sh

cleanup: cleanup-certs ## cleanup example environment
	docker compose -p stevedore-functional-test-manual down -v --remove-orphans

generate-keys: cleanup-keys ## generate an ssh key pair required to autheneticate to git server
	@docker compose -p stevedore-functional-test-manual run --rm openssh -t rsa -q -N "password" -f id_rsa -C "apenella@stevedore.test"

cleanup-keys: ## cleanup the ssh key pair
	@docker compose -p stevedore-functional-test-manual run --rm --entrypoint /bin/sh openssh -c 'rm -rf $$(ls)'

generate-certs: cleanup-certs ## generate certificate for stevedore.test
	@docker compose -p stevedore-functional-test-manual run --rm openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config /root/ssl/stevedore.test.cnf
	@docker compose -p stevedore-functional-test-manual run --rm openssl  x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile /root/ssl/stevedore.test.cnf

cleanup-certs: ## cleanup certificates
	@docker compose -p stevedore-functional-test-manual run --rm --entrypoint /bin/sh openssl -c 'rm -rf $$(ls)'

attach-client: start prepare-images ## attach to client container preloading docker images
	docker compose -p stevedore-functional-test-manual build stevedore
	docker compose -p stevedore-functional-test-manual run -w /app/test/stack/client/stevedore stevedore sh

attach-client-no-images: start ## attach to client container
#attach-client: ## attach to client container
	docker compose -p stevedore-functional-test-manual build stevedore
	docker compose -p stevedore-functional-test-manual run -w /app/test/stack/client/stevedore stevedore sh
