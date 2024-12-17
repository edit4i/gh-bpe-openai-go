# Go Bindings for OpenAI BPE Tokenizer

This package provides Go bindings for the Rust-based OpenAI BPE (Byte Pair Encoding) tokenizer. It supports both the CL100k and O200k tokenizer variants used by OpenAI's models.

## Prerequisites

- Go 1.23 or later
- Rust toolchain (for building the native library)
- cbindgen (for generating C headers)
- Git (for cloning submodules)

## Installation

```bash
# Clone the repository with submodules
git clone --recursive https://github.com/edit4i/gh-bpe-openai-go
cd gh-bpe-openai-go

# If you already cloned without --recursive, run:
git submodule update --init --recursive

# Install cbindgen
cargo install cbindgen

# Build the project
make
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    bpe "github.com/edit4i/gh-bpe-openai-go"
)

func main() {
    // Create a new CL100k tokenizer (used by GPT-4 and ChatGPT)
    tokenizer, err := bpe.NewCL100kTokenizer()
    if err != nil {
        log.Fatal(err)
    }

    text := "Hello, world!"
    
    // Count tokens
    count, err := tokenizer.Count(text)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Token count: %d\n", count)

    // Encode text to tokens
    tokens, err := tokenizer.Encode(text)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Tokens: %v\n", tokens)

    // Decode tokens back to text
    decoded, err := tokenizer.Decode(tokens)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Decoded text: %s\n", decoded)
}
```

## Features

- Fast, native implementation using Rust
- Support for CL100k and O200k tokenizers
- Memory-safe bindings with automatic cleanup
- Token counting with optional limits
- Encode/decode functionality
- Full test coverage

## Project Structure

- `bpe-openai-ffi/`: Rust FFI layer for the tokenizer
- `rust-gems/`: Submodule containing the original Rust implementation
- `bpe.go`: Go bindings
- `bpe_test.go`: Test suite
- `Makefile`: Build automation
- `cbindgen.toml`: C bindings configuration

## License

MIT License - see LICENSE file for details
