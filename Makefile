.PHONY: all clean test bench performance build-all

BUILD_DIR = build
INCLUDE_DIR = include
LIB_DIR = lib

# Platform-specific directories
LINUX_AMD64_DIR = $(LIB_DIR)/linux_amd64
LINUX_ARM64_DIR = $(LIB_DIR)/linux_arm64
DARWIN_AMD64_DIR = $(LIB_DIR)/darwin_amd64
DARWIN_ARM64_DIR = $(LIB_DIR)/darwin_arm64
WINDOWS_AMD64_DIR = $(LIB_DIR)/windows_amd64

all: $(BUILD_DIR)/libbpe_openai_ffi.a $(INCLUDE_DIR)/bpe_openai.h

# Create all platform directories
create-dirs:
	mkdir -p $(LINUX_AMD64_DIR) $(LINUX_ARM64_DIR) $(DARWIN_AMD64_DIR) $(DARWIN_ARM64_DIR) $(WINDOWS_AMD64_DIR)

# Build for all platforms
build-all: create-dirs build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

# Linux AMD64
build-linux-amd64:
	@echo "Building for Linux AMD64..."
	cd bpe-openai-ffi && \
		cargo build --release
	mkdir -p $(LINUX_AMD64_DIR)
	cp bpe-openai-ffi/target/release/libbpe_openai_ffi.so $(LINUX_AMD64_DIR)/

# Linux ARM64
build-linux-arm64:
	@echo "Building for Linux ARM64..."
	cd bpe-openai-ffi && \
		CARGO_TARGET_AARCH64_UNKNOWN_LINUX_GNU_LINKER=aarch64-linux-gnu-gcc \
		cargo build --release --target aarch64-unknown-linux-gnu
	cp bpe-openai-ffi/target/aarch64-unknown-linux-gnu/release/libbpe_openai_ffi.so $(LINUX_ARM64_DIR)/ || \
	cp bpe-openai-ffi/target/release/libbpe_openai_ffi.so $(LINUX_ARM64_DIR)/

# Darwin (macOS) AMD64
build-darwin-amd64:
	@echo "Building for macOS AMD64..."
	cd bpe-openai-ffi && \
		CARGO_TARGET_X86_64_APPLE_DARWIN_LINKER=x86_64-apple-darwin-gcc \
		cargo build --release --target x86_64-apple-darwin
	cp bpe-openai-ffi/target/x86_64-apple-darwin/release/libbpe_openai_ffi.dylib $(DARWIN_AMD64_DIR)/ || \
	cp bpe-openai-ffi/target/release/libbpe_openai_ffi.dylib $(DARWIN_AMD64_DIR)/

# Darwin (macOS) ARM64
build-darwin-arm64:
	@echo "Building for macOS ARM64..."
	cd bpe-openai-ffi && \
		CARGO_TARGET_AARCH64_APPLE_DARWIN_LINKER=aarch64-apple-darwin-gcc \
		cargo build --release --target aarch64-apple-darwin
	cp bpe-openai-ffi/target/aarch64-apple-darwin/release/libbpe_openai_ffi.dylib $(DARWIN_ARM64_DIR)/ || \
	cp bpe-openai-ffi/target/release/libbpe_openai_ffi.dylib $(DARWIN_ARM64_DIR)/

# Windows AMD64
build-windows-amd64:
	@echo "Building for Windows AMD64..."
	cd bpe-openai-ffi && \
		CARGO_TARGET_X86_64_PC_WINDOWS_GNU_LINKER=x86_64-w64-mingw32-gcc \
		cargo build --release --target x86_64-pc-windows-gnu
	cp bpe-openai-ffi/target/x86_64-pc-windows-gnu/release/bpe_openai_ffi.dll $(WINDOWS_AMD64_DIR)/ || \
	cp bpe-openai-ffi/target/release/bpe_openai_ffi.dll $(WINDOWS_AMD64_DIR)/

$(BUILD_DIR)/libbpe_openai_ffi.a: 
	@mkdir -p $(BUILD_DIR)
	cd bpe-openai-ffi && \
		cargo build --release
	cp bpe-openai-ffi/target/release/libbpe_openai_ffi.a $(BUILD_DIR)/

$(INCLUDE_DIR)/bpe_openai.h: $(BUILD_DIR)/libbpe_openai_ffi.a
	@mkdir -p $(INCLUDE_DIR)
	cbindgen --config cbindgen.toml \
		--crate bpe-openai-ffi \
		--output $(INCLUDE_DIR)/bpe_openai.h \
		bpe-openai-ffi

test: build-linux-amd64
	@echo "Running tests with LD_LIBRARY_PATH=$(LINUX_AMD64_DIR)"
	LD_LIBRARY_PATH=$(LINUX_AMD64_DIR) go test -v ./...

bench: build-linux-amd64
	LD_LIBRARY_PATH=$(LINUX_AMD64_DIR) go test -bench=. -benchmem ./...

performance: build-linux-amd64
	@echo "\nDetailed Performance Report:"
	@echo "============================"
	LD_LIBRARY_PATH=$(LINUX_AMD64_DIR) go test -bench=. -benchmem -benchtime=5s -count=5 ./... | tee benchmark.txt
	@echo "\nResults saved to benchmark.txt"

clean:
	rm -rf $(BUILD_DIR) $(INCLUDE_DIR) $(LIB_DIR) benchmark.txt
	cd bpe-openai-ffi && cargo clean

# Install required Rust targets
install-rust-targets:
	rustup target add x86_64-unknown-linux-gnu
	rustup target add aarch64-unknown-linux-gnu
	rustup target add x86_64-apple-darwin
	rustup target add aarch64-apple-darwin
	rustup target add x86_64-pc-windows-gnu

# Install required system dependencies (Ubuntu/Debian)
install-deps-ubuntu:
	sudo apt-get update && sudo apt-get install -y \
		gcc-x86-64-linux-gnu \
		gcc-aarch64-linux-gnu \
		gcc-mingw-w64 \
		osxcross

# Show help
help:
	@echo "Available targets:"
	@echo "  all              - Build the library for the current platform"
	@echo "  build-all        - Build the library for all supported platforms"
	@echo "  build-linux-amd64  - Build for Linux AMD64"
	@echo "  build-linux-arm64  - Build for Linux ARM64"
	@echo "  build-darwin-amd64 - Build for macOS AMD64"
	@echo "  build-darwin-arm64 - Build for macOS ARM64"
	@echo "  build-windows-amd64 - Build for Windows AMD64"
	@echo "  test             - Run tests"
	@echo "  bench            - Run benchmarks"
	@echo "  performance      - Run detailed performance tests"
	@echo "  clean            - Clean build artifacts"
	@echo "  install-rust-targets - Install required Rust targets"
	@echo "  install-deps-ubuntu  - Install required system dependencies (Ubuntu/Debian)"
