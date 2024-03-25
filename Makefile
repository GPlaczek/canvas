MODULE_DIR := $(shell pwd)
BUILD_DIR := $(MODULE_DIR)/build

CANVAS_TARGET := $(BUILD_DIR)/canvas

$(BUILD_DIR):
	mkdir -p $@

.PHONY: build
build: | $(BUILD_DIR)
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
