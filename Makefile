GOCMD=env go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

BINARY=boing
TESTS=./...
COVERAGE_FILE=coverage.out

WASM_DIR=./html/

.PHONY: all test build build-wasm coverage clean resources

all: test build

build:
		$(GOBUILD) -o $(BINARY) -v

build-wasm:
		cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $(WASM_DIR)
		GOOS=js GOARCH=wasm $(GOBUILD) -o $(WASM_DIR)$(BINARY).wasm -v
		gzip --keep --best $(WASM_DIR)$(BINARY).wasm

test:
		$(GOTEST) -v $(TESTS)

coverage:
		$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		$(GOCLEAN)
		rm -f $(BINARY) $(COVERAGE_FILE)
