GO ?= go
CC ?= cc

APP = alt-tab
CMD = ./cmd/alt-tab
BIN_DIR = bin
BUILD_DIR = build
NATIVE_DIR = internal/nativechords
TEST_DIR = tests

TARGET = $(BIN_DIR)/$(APP)
NATIVE_BUILD_DIR = $(BUILD_DIR)/native
NATIVE_TEST_TARGET = $(NATIVE_BUILD_DIR)/test-chord-backend

BUILD ?= debug

GO_SOURCES = $(shell find cmd internal -type f \
             \( -name '*.go' -o -name '*.c' -o -name '*.h' \) 2>/dev/null)
NATIVE_SOURCES = $(NATIVE_DIR)/chord_library.c \
                 $(NATIVE_DIR)/chord_api.c
NATIVE_HEADERS = $(NATIVE_DIR)/chord.h \
                 $(NATIVE_DIR)/chord_library.h \
                 $(NATIVE_DIR)/chord_api.h
NATIVE_OBJECTS = $(NATIVE_SOURCES:$(NATIVE_DIR)/%.c=$(NATIVE_BUILD_DIR)/%.o)
NATIVE_TEST_OBJECT = $(NATIVE_BUILD_DIR)/test_chord_library.o

CFLAGS_BASE = -std=c11 -Wall -Wextra -pedantic
SANITIZERS = -fsanitize=address,undefined
GO_BUILD_FLAGS =

ifeq ($(BUILD),debug)
    CFLAGS = $(CFLAGS_BASE) -g $(SANITIZERS)
    LDFLAGS = $(SANITIZERS)
else ifeq ($(BUILD),release)
    CFLAGS = $(CFLAGS_BASE) -O2 -DNDEBUG
    LDFLAGS =
    GO_BUILD_FLAGS = -trimpath -ldflags="-s -w"
else
    $(error Unknown BUILD type: $(BUILD))
endif

all: build

build: $(TARGET)

check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { \
		echo "error: Go is required but '$(GO)' was not found"; \
		exit 1; \
	}
	@command -v $(firstword $(CC)) >/dev/null 2>&1 || { \
		echo "error: a C compiler is required but '$(firstword $(CC))' was not found"; \
		exit 1; \
	}
	@test "$$($(GO) env CGO_ENABLED)" = "1" || { \
		echo "error: cgo is disabled; run with CGO_ENABLED=1"; \
		exit 1; \
	}
	@$(GO) list -mod=readonly -deps ./... >/dev/null
	@$(GO) mod verify >/dev/null

$(TARGET): check-deps go.mod go.sum $(GO_SOURCES)
	mkdir -p $(dir $@)
	$(GO) build $(GO_BUILD_FLAGS) -o $@ $(CMD)

$(NATIVE_BUILD_DIR)/%.o: $(NATIVE_DIR)/%.c $(NATIVE_HEADERS)
	mkdir -p $(dir $@)
	$(CC) -I$(NATIVE_DIR) $(CFLAGS) -c $< -o $@

$(NATIVE_TEST_OBJECT): $(TEST_DIR)/test_chord_library.c $(NATIVE_HEADERS)
	mkdir -p $(dir $@)
	$(CC) -I$(NATIVE_DIR) $(CFLAGS) -c $< -o $@

$(NATIVE_TEST_TARGET): $(NATIVE_TEST_OBJECT) $(NATIVE_OBJECTS)
	$(CC) $^ $(LDFLAGS) -o $@

test-native: $(NATIVE_TEST_TARGET)
	./$(NATIVE_TEST_TARGET)

test-go: check-deps
	$(GO) test ./...

test: test-native test-go

race: check-deps
	$(GO) test -race ./internal/...

vet: check-deps
	$(GO) vet ./...

run: $(TARGET)
	./$(TARGET)

clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR) $(APP)

.PHONY: all build check-deps clean race run test test-go test-native vet
