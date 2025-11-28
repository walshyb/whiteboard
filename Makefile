# =========================================================================
# Protobuf Settings
# =========================================================================

PROTO_DIR := proto

PROTO_FILES := $(PROTO_DIR)/events.proto

GO_OUT_DIR := $(PROTO_DIR)
GO_OPTS := paths=source_relative

CLIENT_OUT_DIR := client/src/proto/generated
TS_PROTO_PLUGIN := client/node_modules/.bin/protoc-gen-ts_proto
TS_OPTS := esModuleInterop=true,forceLong=string,outputServices=false

.PHONY: all proto clean help

all: proto

clean:
	@echo "Cleaning generated"
	@find $(GO_OUT_DIR) -name "*.pb.go" -delete
	@find $(CLIENT_OUT_DIR) -name "*.ts" -delete

proto: $(PROTO_FILES)
	@echo "Compiling Go Protobuf files"
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT_DIR) \
		--go_opt=$(GO_OPTS) \
		\
		--plugin=$(TS_PROTO_PLUGIN) \
    --ts_proto_opt=$(TS_OPTS) \
		--ts_proto_out=$(CLIENT_OUT_DIR) \
		$(PROTO_FILES)
	@echo "Compiling TS Protobuf files"


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
