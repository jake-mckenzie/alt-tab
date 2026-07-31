CC = gcc
GO = go

SRC_DIR   = src
INC_DIR   = include
BUILD_DIR = build
TEST_DIR  = tests

TARGET = alt-tab
TEST_TARGET = $(BUILD_DIR)/test-alt-tab
TUI_TARGET = $(BUILD_DIR)/alt-tab-tui
TUI_PACKAGE = ./cmd/alt-tab-tui

BUILD ?= debug

APP_SRCS = $(SRC_DIR)/app/main.c
CORE_SRCS = $(SRC_DIR)/theory/chord_library.c \
            $(SRC_DIR)/render/terminal_renderer.c
SRCS = $(APP_SRCS) $(CORE_SRCS)
TEST_SRCS = $(TEST_DIR)/test_chord_library.c
GO_SRCS = $(shell find cmd internal -name '*.go' 2>/dev/null)

OBJS = $(SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
CORE_OBJS = $(CORE_SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
TEST_OBJS = $(TEST_SRCS:$(TEST_DIR)/%.c=$(BUILD_DIR)/tests/%.o)

CFLAGS_BASE = -std=c11 -Wall -Wextra -pedantic
CPPFLAGS = -I$(INC_DIR)
SANITIZERS = -fsanitize=address,undefined

ifeq ($(BUILD),debug)
    CFLAGS  = $(CFLAGS_BASE) -g $(SANITIZERS)
    LDFLAGS = $(SANITIZERS)
else ifeq ($(BUILD),release)
    CFLAGS  = $(CFLAGS_BASE) -O2 -DNDEBUG
    LDFLAGS =
else
    $(error Unknown BUILD type: $(BUILD))
endif

all: $(TARGET)

tui: $(TUI_TARGET)

$(TARGET): $(OBJS)
	$(CC) $(OBJS) $(LDFLAGS) -o $@

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c
	mkdir -p $(dir $@)
	$(CC) $(CPPFLAGS) $(CFLAGS) -c $< -o $@

$(BUILD_DIR)/tests/%.o: $(TEST_DIR)/%.c
	mkdir -p $(dir $@)
	$(CC) $(CPPFLAGS) $(CFLAGS) -c $< -o $@

$(TEST_TARGET): $(TEST_OBJS) $(CORE_OBJS)
	$(CC) $(TEST_OBJS) $(CORE_OBJS) $(LDFLAGS) -o $@

$(TUI_TARGET): go.mod go.sum $(GO_SRCS)
	mkdir -p $(dir $@)
	$(GO) build -o $@ $(TUI_PACKAGE)

test: $(TARGET) $(TEST_TARGET)
	./$(TEST_TARGET)
	printf 'q\n' | ./$(TARGET) --full-neck | grep -q 'Mode    : full neck'
	printf 'q\n' | ./$(TARGET) -f | grep -q 'Mode    : full neck'
	printf '%s\n' '-f' 'q' | ./$(TARGET) | grep -q 'Display changed to full neck'
	printf '%s\n' '--full-neck' 'q' | ./$(TARGET) | grep -q 'Display changed to full neck'
	./$(TARGET) A -f | grep -q ' 27 '
	./$(TARGET) -f A | grep -q ' 27 '
	printf '%s\n' 'A -f' 'q' | ./$(TARGET) | grep -q ' 27 '
	./$(TARGET) C:2 | grep -q 'Chord: C (variation 2)'
	./$(TARGET) C -a | grep -q 'Chord: C (variation 1)'
	./$(TARGET) C -a | grep -q 'Chord: C (variation 2)'
	printf '%s\n' 'C:2' 'q' | ./$(TARGET) | grep -q 'Chord: C (variation 2)'
	./$(TARGET) --help | grep -q 'SUPPORTED CHORDS'
	$(GO) test ./...

clean:
	rm -rf $(BUILD_DIR) $(TARGET)

run: $(TARGET)
	./$(TARGET) $(ARGS)

run-full: $(TARGET)
	./$(TARGET) --full-neck

run-tui: $(TUI_TARGET)
	./$(TUI_TARGET)

.PHONY: all clean run run-full run-tui test tui
