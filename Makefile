GO ?= go

APP = alt-tab
CMD = ./cmd/alt-tab
BIN_DIR = bin

TARGET = $(BIN_DIR)/$(APP)

BUILD ?= debug

GO_FILES = $(shell find cmd internal -type f -name '*.go' 2>/dev/null)
GO_BUILD_FLAGS =

ifeq ($(BUILD),release)
    GO_BUILD_FLAGS = -trimpath -ldflags="-s -w"
else ifneq ($(BUILD),debug)
    $(error Unknown BUILD type: $(BUILD))
endif

all: build

build: check-deps go.mod go.sum $(GO_FILES)
	mkdir -p $(dir $(TARGET))
	$(GO) build $(GO_BUILD_FLAGS) -o $(TARGET) $(CMD)

fmt:
	$(GO) fmt ./...

fmt-check:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	test -z "$$unformatted" || { \
		echo "error: these Go files require formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	}

check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { \
		echo "error: Go is required but '$(GO)' was not found"; \
		exit 1; \
	}
	@$(GO) list -mod=readonly -deps ./... >/dev/null
	@$(GO) mod verify >/dev/null

test-go: check-deps
	$(GO) test ./...

test: test-go

race: check-deps
	$(GO) test -race ./internal/...

vet: check-deps
	$(GO) vet ./...

check: fmt-check test race vet

run: build
	./$(TARGET)

clean:
	rm -rf $(BIN_DIR)

.PHONY: all build check check-deps clean fmt fmt-check race run test test-go vet
