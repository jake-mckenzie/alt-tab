# Compiler (Apple Clang on macOS)
CC = gcc

# Directories
SRC_DIR   = src
INC_DIR   = include
BUILD_DIR = build

# Target name
TARGET = alt-tab

# Build type: debug (default) or release
BUILD ?= debug

# Source files
SRCS = $(SRC_DIR)/main.c \
       $(SRC_DIR)/matrix.c \
       $(SRC_DIR)/chords.c

# Object files
OBJS = $(SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)

# Common compiler flags
CFLAGS_BASE = -std=c11 -Wall -Wextra -pedantic -I$(INC_DIR)

# Sanitizers (debug only)
SANITIZERS = -fsanitize=address,undefined

# Build configuration
ifeq ($(BUILD),debug)
    CFLAGS  = $(CFLAGS_BASE) -g
    LDFLAGS = $(SANITIZERS)
else ifeq ($(BUILD),release)
    CFLAGS  = $(CFLAGS_BASE) -O2 -DNDEBUG
    LDFLAGS =
else
    $(error Unknown BUILD type: $(BUILD))
endif

# Default target
all: $(TARGET)

# Link step (IMPORTANT: LDFLAGS included here)
$(TARGET): $(OBJS)
	$(CC) $(OBJS) $(LDFLAGS) -o $@

# Compile step
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c | $(BUILD_DIR)
	$(CC) $(CFLAGS) -c $< -o $@

# Create build directory if needed
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR) $(TARGET)

# Run program
run: $(TARGET)
	./$(TARGET)

.PHONY: all clean run
