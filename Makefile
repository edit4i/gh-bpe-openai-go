.PHONY: all clean test bench performance

BUILD_DIR = build
INCLUDE_DIR = include

all: $(BUILD_DIR)/libbpe_openai_ffi.a $(INCLUDE_DIR)/bpe_openai.h

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

test: all
	go test -v ./...

bench: all
	go test -bench=. -benchmem ./...

performance: bench
	@echo "\nDetailed Performance Report:"
	@echo "============================"
	go test -bench=. -benchmem -benchtime=5s -count=5 ./... | tee benchmark.txt
	@echo "\nResults saved to benchmark.txt"

clean:
	rm -rf $(BUILD_DIR) $(INCLUDE_DIR) benchmark.txt
	cd bpe-openai-ffi && cargo clean
