# Root path for shared proto files
PROTO_DIR = shared_proto

# Target directories for generated proto files
API_GATEWAY_PROTO = api_gateway/proto

# List of services to generate
SERVICES = auth_service user_service author_service category_service book_service

# Default target: Display help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make generate-all             Generate proto files for all services."
	@echo "  make generate-proto SERVICE=<service_name> Generate proto files for a specific service."
	@echo ""
	@echo "Available services: $(SERVICES)"

# Generate proto files for all services
.PHONY: generate-all
generate-all:
	@echo "Generating proto files for all services..."
	@for SERVICE in $(SERVICES); do \
		$(MAKE) generate-proto SERVICE=$$SERVICE; \
	done
	@echo "All proto files have been generated."

# Generate proto files for a specific service
.PHONY: generate-proto
generate-proto:
	@echo "Generating proto files for $(SERVICE)..."
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE flag is required. Usage: make generate-proto SERVICE=<service_name>"; \
		exit 1; \
	fi
	@if ! echo "$(SERVICES)" | grep -q -w "$(SERVICE)"; then \
		echo "Error: Unknown service '$(SERVICE)'. Available services: $(SERVICES)"; \
		exit 1; \
	fi

	@# Create target directories if they don't exist
	@mkdir -p $(API_GATEWAY_PROTO)/$(SERVICE)
	@mkdir -p $(SERVICE)/proto

	@# Generate proto files for both API Gateway and Service Directory
	@for TARGET_DIR in $(API_GATEWAY_PROTO)/$(SERVICE) $(SERVICE)/proto; do \
		echo "Generating files in $$TARGET_DIR..."; \
		protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$$TARGET_DIR \
		--go_opt=paths=source_relative \
		--go-grpc_out=$$TARGET_DIR \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/$(SERVICE).proto; \
	done

	@echo "Proto files generated for $(SERVICE)."
