# this tells Make to run 'make help' if the user runs 'make'
# without this, Make would use the first target as the default
.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


createSudoersFile: ## Creates the sudoers file
	@echo "Creating sudoers file"
	@source ./common.sh && createSudoersFile

installCloudProviderKind: ## Install cloud-provider-kind
	@echo "Install cloud-provider-kind"
	@source ./common.sh && installCloudProviderKind

build: ## Builds the tray application
	@echo "Building tray application"
	@source ./common.sh && installCloudProviderKindTray

install: installCloudProviderKind createSudoersFile build ## Installs the tray application

run: ## Runs the tray application
	@echo "Running tray application"
	@cloud-provider-kind-tray
