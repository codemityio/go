.PHONY: cmd docs tools

-include Makefile.env
-include Makefile.env.local

%: # Silence errors about non existing targets
	@true

default: # A default target to initiate interactive menu
	@source scripts/func.sh && check "git docker go goimports gofumpt"
	@scripts/default.sh make

help: ## Prints help for targets with comments
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

version: ## Print the most recent version
	@scripts/tools.sh version

next: ## Create a new version (bump prerelease or patch)
	@scripts/tools.sh next

check: prep gen fmt statan test-race cov diff cleanup ## Run all CI required targets

###########
## D E V ##
###########

prep: ## Prepare dev tools
	@scripts/tools.sh prep

cmd: ## Run a command passed as COMMAND= value (e.g. make cmd COMMAND="make check")
	@scripts/tools.sh cmd

cleanup: ## Cleanup project
	@scripts/tools.sh cleanup

update: ## Update all dependencies
	@scripts/tools.sh update

vendor: ## Run go mod vendor
	@go mod vendor -v

gen: ## Go generate
	@scripts/tools.sh gen

fmt: ## Format code
	@scripts/tools.sh fmt

statan: ## Analyze code
	@scripts/tools.sh statan

statan-fix: ## Analyze code and fix
	@scripts/tools.sh statan-fix

test: ## Run tests
	@scripts/tools.sh test

test-race: ## Run race tests
	@scripts/tools.sh test-race

cov: ## Check coverage
	@scripts/tools.sh cov

cov-report: ## Check coverage report
	@scripts/tools.sh cov-report

cov-open: ## Inspect coverage in the browser
	@scripts/tools.sh cov-open

diff: ## Check diff to ensure this project consistency
	@scripts/tools.sh diff

#############
## D O C S ##
#############

docs: docs-main ## Generate all docs
	@PACKAGES='$(shell find "${PWD}/pkg" -mindepth 1 -maxdepth 1 -type d -exec basename {} \; 2>/dev/null)' make docs-uml docs-depgraph docs-pkg docs-render

docs-uml: ## Generate UML documentation
	@scripts/docs.sh uml

docs-depgraph: ## Generate dependency graph
	@scripts/docs.sh depgraph

docs-pkg: ## Generate pkg docs
	@scripts/docs.sh pkg

docs-render: ## Render diagrams
	@scripts/docs.sh render

docs-main: ## Generate main docs
	@scripts/docs.sh main

##########################
## D A N G E R  Z O N E ##
##########################

reset: ## Stop and remove project containers, remove project volumes, remove project images
	@docker ps -a --filter "name=${BASE_NAME}" --format "{{.ID}}" | xargs -r docker stop
	@docker ps -a --filter "name=${BASE_NAME}" --format "{{.ID}}" | xargs -r docker rm
	@docker volume ls --filter "name=${BASE_NAME}" --format "{{.Name}}" | xargs -r docker volume rm
	@docker system prune -f
