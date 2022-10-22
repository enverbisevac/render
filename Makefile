ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all
all: help

# Format go code and error if any changes are made
.PHONY: format
format: $(GOPATH)/bin/goimports ## Format files using goimports
	@echo "Runing gofumpt"
	@gofumpt -l -w .
	@echo "Running goimports"
	@test -z $$(goimports -w ./..) || (echo "goimports would make a change. Please verify and commit the proposed changes"; exit 1)

.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint ## Run golangci-lint
	@golangci-lint run -v

.PHONY: test
test: ## Run unit tests with coverage
	@go test -v -race -coverprofile=coverage.out -covermode=atomic

.PHONY: coverage
coverage: test
	@go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: test
	@go tool cover -html=coverage.out

##################################################
# tools managment rules
##################################################
$(GOPATH)/bin/golangci-lint:
	@echo "🔘 Installing golangci-lint... (`date '+%H:%M:%S'`)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

$(GOPATH)/bin/goimports:
	@echo "🔘 Installing goimports ... (`date '+%H:%M:%S'`)"
	@go install golang.org/x/tools/cmd/goimports@latest

$(GOPATH)/bin/gofumpt:
	@echo "🔘 Installing gofumpt ... (`date '+%H:%M:%S'`)"
	@go install mvdan.cc/gofumpt@latest

.PHONY: tools
tools: $(GOPATH)/bin/golangci-lint $(GOPATH)/bin/goimports $(GOPATH)/bin/gofumpt ## Install tools

.PHONY: update-tools
update-tools: delete-tools tools ## Update tools

.PHONY: delete-tools
delete-tools: ## Delete installed tools
	@rm $(GOPATH)/bin/golangci-lint || true
	@rm $(GOPATH)/bin/goimports || true
	@rm $(GOPATH)/bin/gofumpt || true

## Help:
.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)