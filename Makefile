CC = gcc

SRC_DIR   = src
INC_DIR   = include
BUILD_DIR = build
TEST_DIR  = tests

TARGET = alt-tab
TEST_TARGET = $(BUILD_DIR)/test-alt-tab

BUILD ?= debug

APP_SRCS = $(SRC_DIR)/app/main.c
CORE_SRCS = $(SRC_DIR)/theory/chord_library.c \
            $(SRC_DIR)/render/terminal_renderer.c
SRCS = $(APP_SRCS) $(CORE_SRCS)
TEST_SRCS = $(TEST_DIR)/test_chord_library.c

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

test: $(TEST_TARGET)
	./$(TEST_TARGET)

clean:
	rm -rf $(BUILD_DIR) $(TARGET)

run: $(TARGET)
	./$(TARGET)

.PHONY: all clean run test
