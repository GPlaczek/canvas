MODULE_DIR := $(shell pwd)
BUILD_DIR := $(MODULE_DIR)/build

CANVAS_TARGET := $(BUILD_DIR)/canvas

$(BUILD_DIR):
	mkdir -p $@

$(MODULE_DIR)/front/messages.proto: $(MODULE_DIR)/proto/messages.proto
	cp -- $< $@

$(MODULE_DIR)/pkg/protocol/messages.pb.go: ./proto/messages.proto
	protoc -I=proto --go_out=pkg $<

.PHONY: build
build: $(MODULE_DIR)/pkg/protocol/messages.pb.go | $(BUILD_DIR)
	go build -o $(BUILD_DIR)/canvas $(MODULE_DIR)/cmd/canvas.go

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: format
format:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: front
front: $(MODULE_DIR)/front/messages.proto
