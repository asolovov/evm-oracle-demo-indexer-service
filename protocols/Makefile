.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Installation targets
.PHONY: buf-install
buf-install: ## Install Buf CLI
	@echo "Installing Buf..."
	@which buf > /dev/null || go install github.com/bufbuild/buf/cmd/buf@latest
	@echo "Buf installed successfully"

.PHONY: protoc-install
protoc-install: ## Install protoc plugins
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Protoc plugins installed successfully"

.PHONY: install
install: buf-install protoc-install ## Install all required tools

# Buf workflow targets
.PHONY: buf-lint
buf-lint: ## Run Buf linting on all modules
	@echo "Running Buf lint..."
	@buf lint

.PHONY: buf-breaking
buf-breaking: ## Check for breaking changes against main branch
	@echo "Checking for breaking changes..."
	@if git rev-parse --verify origin/main >/dev/null 2>&1; then \
		buf breaking --against '.git#branch=origin/main'; \
	else \
		echo "No main branch found, skipping breaking change detection"; \
	fi

.PHONY: buf-format
buf-format: ## Format proto files using Buf
	@echo "Formatting proto files..."
	@buf format -w

.PHONY: buf-generate
buf-generate: ## Generate Go code using Buf (requires PACKAGE)
ifndef PACKAGE
	@echo "Error: PACKAGE parameter is required"
	@echo "Usage: make buf-generate PACKAGE=price"
	@exit 1
endif
	@echo "Generating Go code for $(PACKAGE) using Buf..."
	@buf generate --path $(PACKAGE)

.PHONY: buf-generate-all
buf-generate-all: ## Generate Go code for all packages using Buf
	@echo "Generating Go code for all packages using Buf..."
	@buf generate

# Protoc workflow targets
.PHONY: protoc-validate
protoc-validate: ## Validate proto files using protoc
ifdef PACKAGE
	@echo "Validating $(PACKAGE) using protoc..."
	@find $(PACKAGE) -name "*.proto" -exec protoc --proto_path=. --proto_path=$(PACKAGE) --descriptor_set_out=/dev/null {} \;
else
	@echo "Validating all proto files using protoc..."
	@find . -name "*.proto" -exec protoc --proto_path=. --descriptor_set_out=/dev/null {} \;
endif

.PHONY: protoc-generate
protoc-generate: ## Generate Go code using protoc (requires PACKAGE)
ifndef PACKAGE
	@echo "Error: PACKAGE parameter is required"
	@echo "Usage: make protoc-generate PACKAGE=price"
	@exit 1
endif
	@echo "Generating Go code for $(PACKAGE) using protoc..."
	@find $(PACKAGE) -name "*.proto" -exec protoc \
		--proto_path=. \
		--proto_path=$(PACKAGE) \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		{} \;

.PHONY: protoc-generate-all
protoc-generate-all: ## Generate Go code for all packages using protoc
	@echo "Generating Go code for all packages using protoc..."
	@for dir in common price oracle indexer; do \
		echo "Generating $$dir..."; \
		$(MAKE) protoc-generate PACKAGE=$$dir; \
	done

# Validation targets
.PHONY: validate
validate: buf-lint protoc-validate ## Run all validation checks

.PHONY: validate-ci
validate-ci: buf-lint buf-breaking protoc-validate ## Run all validation checks for CI

# Service scaffolding
.PHONY: add-service
add-service: ## Scaffold a new service (requires NAME)
ifndef NAME
	@echo "Error: NAME parameter is required"
	@echo "Usage: make add-service NAME=product"
	@exit 1
endif
	@./scripts/add-service.sh $(NAME)

# Clean targets
.PHONY: clean
clean: ## Remove all generated files
	@echo "Removing generated files..."
	@find . -name "*.pb.go" -type f -delete
	@find . -name "*_grpc.pb.go" -type f -delete
	@echo "Clean complete"

.PHONY: clean-buf
clean-buf: ## Remove Buf cache
	@echo "Removing Buf cache..."
	@rm -rf .buf/
	@echo "Buf cache removed"

# Default target
.PHONY: all
all: validate ## Run lint and validation (default)

# Generate targets (shortcuts)
.PHONY: generate
generate: buf-generate-all ## Generate all code using Buf (default generator)

.PHONY: gen
gen: generate ## Alias for generate
