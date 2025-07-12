.PHONY=all build run


BINARY_NAME=echo
BUILD_PATH=cmd/$(BINARY_NAME)
BINARY_PATH=bin/$(BINARY_NAME)

all: build
	run

build:
	@echo "Building project..."
	@go build -o $(BINARY_PATH) ./$(BUILD_PATH)

run: build
	@echo "Running binary..."
	@./$(BINARY_PATH)

clean:
	@echo "Deleting binaries..."
	@rm bin/*
