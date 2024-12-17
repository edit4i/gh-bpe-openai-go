# bpe-openai-go

Go bindings for OpenAI's BPE tokenizer implemented in Rust. This project provides high-performance tokenization capabilities for Go applications, offering significant speed improvements over pure Go implementations.

This project wraps GitHub's [Rust BPE implementation](https://github.com/github/rust-gems/tree/main/crates/bpe), which is part of their `rust-gems` collection. The implementation is known for its speed, correctness, and novel algorithms for Byte Pair Encoding.

## Performance

Our Rust-based implementation shows significant performance improvements compared to the pure Go implementation ([tiktoken-go](https://github.com/pkoukk/tiktoken-go)):

![Benchmark Results](benchmark/results.png)

Key findings:
- 4-6x faster than tiktoken-go across all text types
- CL100k model shows ~393-540% improvement
- O200k model shows even better results with ~483-626% improvement
- Consistent performance gains for both short and long texts
- Particularly efficient with Unicode text, showing a 622% improvement for O200k

## Installation

```bash
go get github.com/edit4i/gh-bpe-openai-go
```

## Usage

Here's a simple example of how to use the tokenizer:

```go
package main

import (
    "fmt"
    "log"
    
    bpe "github.com/edit4i/gh-bpe-openai-go"
)

func main() {
    tokenizer, err := bpe.NewCL100kTokenizer()
    if err != nil {
        log.Fatal(err)
    }
    defer tokenizer.Free()

    text := "Hello, world!"
    tokens, err := tokenizer.Encode(text)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tokens: %v\n", tokens)
}
```

Check out the [examples](examples/) directory for more usage examples.

## Available Models

- `CL100k` - Used by GPT-4, GPT-3.5-turbo, text-embedding-3-*
- `O200k` - Used by earlier models like GPT-3

## Project Structure

The project is organized as follows:
- `bpe-openai-ffi/` - Rust FFI layer that provides C bindings
- `rust-gems/` - Submodule containing GitHub's Rust BPE implementation
- `benchmark/` - Performance benchmarks and visualization
- `examples/` - Usage examples
- Root files:
  - `bpe.go` - Main Go bindings
  - `bpe_test.go` - Test suite
  - `cbindgen.toml` - C bindings configuration

## Development

### Prerequisites

- Go 1.21+
- Rust (latest stable)
- cbindgen (for generating C headers)

### Building

```bash
# Clone with submodules
git clone --recursive git@github.com:edit4i/gh-bpe-openai-go.git
cd gh-bpe-openai-go

# If you already cloned without --recursive:
git submodule update --init --recursive

# Build
make build
```

### Running Tests

```bash
make test
```

### Running Benchmarks

```bash
# Run the benchmarks
go test -bench=.

# Generate benchmark visualization (requires Python)
make performance # will generate "benchmark.txt"
mv benchmark.txt benchmark/benchmark.txt

cd benchmark
pip install -r requirements.txt
python plot_benchmark.py
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

This project wouldn't be possible without:
- [GitHub's Rust Gems](https://github.com/github/rust-gems) - The underlying Rust implementation
- [tiktoken-go](https://github.com/pkoukk/tiktoken-go) - For benchmark comparisons and inspiration
