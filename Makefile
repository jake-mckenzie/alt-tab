CC = gcc
GO = go

SRC_DIR = src
INC_DIR = include
BUILD_DIR = build
TEST_DIR = tests

TARGET = alt-tab
TEST_TARGET = $(BUILD_DIR)/test-chord-backend
GO_PACKAGE = ./cmd/alt-tab

BUILD ?= debug

BACKEND_SRCS = $(SRC_DIR)/theory/chord_library.c \
               $(SRC_DIR)/backend/chord_api.c
TEST_SRCS = $(TEST_DIR)/test_chord_library.c
GO_SRCS = $(shell find cmd internal -type f \
          \( -name '*.go' -o -name '*.c' \) 2>/dev/null)

BACKEND_OBJS = $(BACKEND_SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
TEST_OBJS = $(TEST_SRCS:$(TEST_DIR)/%.c=$(BUILD_DIR)/tests/%.o)

CFLAGS_BASE = -std=c11 -Wall -Wextra -pedantic
CPPFLAGS = -I$(INC_DIR)
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

all: $(TARGET)

$(TARGET): FORCE go.mod go.sum $(GO_SRCS) $(BACKEND_SRCS)
	$(GO) build $(GO_BUILD_FLAGS) -o $@ $(GO_PACKAGE)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c
	mkdir -p $(dir $@)
	$(CC) $(CPPFLAGS) $(CFLAGS) -c $< -o $@

$(BUILD_DIR)/tests/%.o: $(TEST_DIR)/%.c
	mkdir -p $(dir $@)
	$(CC) $(CPPFLAGS) $(CFLAGS) -c $< -o $@

$(TEST_TARGET): $(TEST_OBJS) $(BACKEND_OBJS)
	$(CC) $(TEST_OBJS) $(BACKEND_OBJS) $(LDFLAGS) -o $@

test: $(TEST_TARGET)
	./$(TEST_TARGET)
	$(GO) test ./...

clean:
	rm -rf $(BUILD_DIR) $(TARGET)

run: $(TARGET)
	./$(TARGET)

FORCE:

.PHONY: FORCE all clean run test
