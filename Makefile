GO ?= go

APP = alt-tab
CMD = ./cmd/alt-tab
RAYLIB_APP = alt-tab-raylib
RAYLIB_CMD = ./cmd/alt-tab-raylib
BIN_DIR = bin

TARGET = $(BIN_DIR)/$(APP)
RAYLIB_TARGET = $(BIN_DIR)/$(RAYLIB_APP)
RAYLIB_TAGS = raylib,x11

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

build-raylib: check-raylib-deps go.mod go.sum $(GO_FILES)
	mkdir -p $(dir $(RAYLIB_TARGET))
	$(GO) build -tags $(RAYLIB_TAGS) $(GO_BUILD_FLAGS) -o $(RAYLIB_TARGET) $(RAYLIB_CMD)

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

check-raylib-deps: check-deps
	@test "$$($(GO) env CGO_ENABLED)" = "1" || { \
		echo "error: the Raylib build requires cgo (set CGO_ENABLED=1)"; \
		exit 1; \
	}
	@command -v "$$($(GO) env CC)" >/dev/null 2>&1 || { \
		echo "error: the Raylib build requires the C compiler reported by 'go env CC'"; \
		exit 1; \
	}
	@$(GO) list -tags $(RAYLIB_TAGS) -mod=readonly -deps $(RAYLIB_CMD) >/dev/null

test-go: check-deps
	$(GO) test ./...

test: test-go

test-raylib: check-raylib-deps
	$(GO) test -tags $(RAYLIB_TAGS) ./internal/rayui

race: check-deps
	$(GO) test -race ./internal/...

vet: check-deps
	$(GO) vet ./...

check: fmt-check test test-raylib race vet

run: build
	./$(TARGET)

run-raylib: build-raylib
	./$(RAYLIB_TARGET)

clean:
	rm -rf $(BIN_DIR)

.PHONY: all build build-raylib check check-deps check-raylib-deps clean fmt fmt-check race run run-raylib test test-go test-raylib vet
