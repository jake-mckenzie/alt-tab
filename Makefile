CC ?= cc
CFLAGS ?= -std=c11 -Wall -Wextra -Wpedantic -Werror -O2
CPPFLAGS = -Iinclude
RAYLIB_CFLAGS = $(shell pkg-config --cflags raylib)
RAYLIB_LIBS = $(shell pkg-config --libs raylib)

APP = alt-tab-raylib
BIN_DIR = bin
TARGET = $(BIN_DIR)/$(APP)
TEST_TARGET = $(BIN_DIR)/test-chords
C_SOURCES = $(shell find src -type f -name '*.c')

all: build

build: check-raylib-deps $(TARGET)

$(TARGET): $(C_SOURCES) include/alt_tab.h
	mkdir -p $(BIN_DIR)
	$(CC) $(CFLAGS) $(CPPFLAGS) $(RAYLIB_CFLAGS) $(C_SOURCES) -o $@ $(RAYLIB_LIBS) -lm

$(TEST_TARGET): tests/test_chords.c src/chords.c include/alt_tab.h
	mkdir -p $(BIN_DIR)
	$(CC) $(CFLAGS) $(CPPFLAGS) tests/test_chords.c src/chords.c -o $@

test: $(TEST_TARGET)
	$(TEST_TARGET)

check: test

check-raylib-deps:
	@command -v $(CC) >/dev/null 2>&1 || { echo "error: a C compiler is required"; exit 1; }
	@pkg-config --exists raylib || { echo "error: raylib development files are required"; exit 1; }

run: build
	./$(TARGET)

clean:
	rm -rf $(BIN_DIR)

.PHONY: all build check check-raylib-deps clean run test
