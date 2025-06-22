.PHONY: init network env up down restart logs \
        registry-init registry-install registry-up registry-down registry-logs \
        git-up git-down git-logs \
        forgejo-setup forgejo-template forgejo-init

# Network config
NETWORK_NAME=git-net
SUBNET=172.30.0.0/16
GATEWAY=172.30.0.1

# Env
ENV_FILE=.env
ENV_INTERNAL_FILE=.env.internal
ENV_FILES=--env-file $(ENV_FILE) --env-file $(ENV_INTERNAL_FILE)

# Directories
REGISTRY_DIR=./registry
GIT_DIR=./git

### --- Shared Top-Level Commands ---

init:
	@echo "🔥 Removing all docker volumes and networks..."
	-docker volume rm $$(docker volume ls -q) || true
	docker network prune -f

network:
	@echo "🧱 Creating Docker network"
	-docker network rm $(NETWORK_NAME)
	docker network create --driver bridge --subnet $(SUBNET) --gateway $(GATEWAY) $(NETWORK_NAME)

env:
	@echo "✍️ Writing .env.internal with HOST_EXTERNAL_IP"
	@echo "HOST_INTERNAL_IP=$(GATEWAY)" > $(ENV_INTERNAL_FILE)
	@echo "HOST_EXTERNAL_IP=$$(ifconfig en0 | grep 'inet ' | awk '{print $$2}')" >> $(ENV_INTERNAL_FILE)

up: network env
	@echo "🚀 Starting Registry and Git"
	make registry-up
	make git-up

down:
	@echo "⬇️ Stopping Registry and Git"
	make registry-down
	make git-down

restart: down env
	@echo "🔄 Restarting all services"
	make registry-up
	make git-up

logs:
	@echo "📜 Tailing all logs"
	make registry-logs
	make git-logs

### --- Registry Commands ---

registry-init:
	@echo "📝 harbor.yml templating ..."
	docker run --rm \
		$(ENV_FILES) \
		-v $(CURDIR)/$(REGISTRY_DIR):/workdir \
		alpine \
		sh -c "apk add --no-cache gettext && envsubst < /workdir/harbor.yml.tmpl > /workdir/harbor.yml"

registry-install:
	cd $(REGISTRY_DIR) && ./install.sh

registry-up:
	docker compose $(ENV_FILES) -f $(REGISTRY_DIR)/docker-compose.yml up -d

registry-down:
	docker compose -f $(REGISTRY_DIR)/docker-compose.yml down

registry-logs:
	docker compose -f $(REGISTRY_DIR)/docker-compose.yml logs -f

### --- GIT Commands ---

git-up:
	docker compose $(ENV_FILES) -f $(GIT_DIR)/docker-compose.yml up -d

git-down:
	docker compose -f $(GIT_DIR)/docker-compose.yml down

git-logs:
	docker compose -f $(GIT_DIR)/docker-compose.yml logs -f

### --- Forgejo Setup Commands ---

forgejo-setup:
	@echo "🚀 Running Forgejo for manual setup (Ctrl+C to exit)"
	docker run --rm -it \
		-p 3000:3000 \
		-p 222:22 \
		-v $(CURDIR)/git/data/forgejo:/data \
		codeberg.org/forgejo/forgejo:11.0.0

forgejo-template:
	@echo "📝 Creating Forgejo app.ini.template"
	mkdir -p ./git/data/forgejo/gitea/conf
	cp ./git/conf/forgejo/app.ini.template ./git/data/forgejo/gitea/conf/app.ini.template

forgejo-init:
	@echo "🔧 Initializing Forgejo app.ini from template"
	docker run --rm \
		$(ENV_FILES) \
		-v $(CURDIR)/git/data/forgejo:/data \
		alpine \
		sh -c "apk add --no-cache gettext && envsubst < /data/gitea/conf/app.ini.template > /data/gitea/conf/app.ini"