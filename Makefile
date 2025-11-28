# =========================================================================
# Protobuf Settings
# =========================================================================

PROTO_DIR := proto

PROTO_FILES := $(PROTO_DIR)/events.proto

GO_OUT_DIR := $(PROTO_DIR)

.PHONY: all proto clean

all: proto

clean:
	@echo "Cleaning generated Go files..."
	@find $(PROTO_DIR) -name "*.pb.go" -delete

proto: $(PROTO_FILES)
	@echo "Compiling Protobuf files to Go..."
	protoc --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		   --proto_path=$(PROTO_DIR) \
		   $(PROTO_FILES)
	@echo "Protobuf compilation complete."

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
