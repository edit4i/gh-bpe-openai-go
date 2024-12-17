// Package bpe provides high-performance Go bindings for OpenAI's BPE (Byte-Pair Encoding) tokenizer.
//
// This package offers a Rust-based implementation that is 4-6x faster than pure Go implementations.
// It supports both CL100k and O200k models, with full Unicode support and efficient memory management.
//
// Example usage:
//
//	tokenizer, err := bpe.NewCL100kTokenizer()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Count tokens
//	count, _ := tokenizer.Count("Hello 👋 World!")
//	fmt.Printf("Token count: %d\n", count)
//
//	// Encode text to tokens
//	tokens, _ := tokenizer.Encode("Hello 👋 World!")
//	fmt.Printf("Tokens: %v\n", tokens)
//
//	// Decode tokens back to text
//	text, _ := tokenizer.Decode(tokens)
//	fmt.Printf("Text: %s\n", text)
package bpe

/*
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -lbpe_openai_ffi
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/lib/linux_arm64 -lbpe_openai_ffi
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/lib/darwin_amd64 -lbpe_openai_ffi
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lbpe_openai_ffi
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -lbpe_openai_ffi
#cgo CFLAGS: -I${SRCDIR}/include
#include <bpe_openai.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// Tokenizer represents a BPE tokenizer instance that provides methods for encoding,
// decoding, and counting tokens in text. The tokenizer automatically manages its
// own memory through Go's garbage collection system.
//
// Use NewCL100kTokenizer() or NewO200kTokenizer() to create a new instance.
// The tokenizer is safe for concurrent use across multiple goroutines.
type Tokenizer struct {
	ptr *C.struct_bpe_TokenizerHandle
}

// Error definitions for common tokenizer operations
var (
	// ErrInvalidTokenizer is returned when attempting to use an uninitialized or freed tokenizer
	ErrInvalidTokenizer = errors.New("invalid tokenizer")
	// ErrEncoding is returned when the tokenizer fails to encode text into tokens
	ErrEncoding = errors.New("encoding error")
	// ErrDecoding is returned when the tokenizer fails to decode tokens back into text
	ErrDecoding = errors.New("decoding error")
)

// NewCL100kTokenizer creates a new CL100k tokenizer instance.
// CL100k is OpenAI's GPT-4 tokenizer, suitable for most modern GPT models.
//
// The tokenizer is automatically freed when it's no longer referenced
// and garbage collected.
//
// Returns:
//   - (*Tokenizer, nil) on success
//   - (nil, ErrInvalidTokenizer) if initialization fails
func NewCL100kTokenizer() (*Tokenizer, error) {
	ptr := C.bpe_cl100k_base()
	if ptr == nil {
		return nil, ErrInvalidTokenizer
	}
	t := &Tokenizer{ptr: (*C.struct_bpe_TokenizerHandle)(ptr)}
	runtime.SetFinalizer(t, (*Tokenizer).free)
	return t, nil
}

// NewO200kTokenizer creates a new O200k tokenizer instance.
// O200k is a newer tokenizer model that provides better handling of
// Unicode text and special characters.
//
// The tokenizer is automatically freed when it's no longer referenced
// and garbage collected.
//
// Returns:
//   - (*Tokenizer, nil) on success
//   - (nil, ErrInvalidTokenizer) if initialization fails
func NewO200kTokenizer() (*Tokenizer, error) {
	ptr := C.bpe_o200k_base()
	if ptr == nil {
		return nil, ErrInvalidTokenizer
	}
	t := &Tokenizer{ptr: (*C.struct_bpe_TokenizerHandle)(ptr)}
	runtime.SetFinalizer(t, (*Tokenizer).free)
	return t, nil
}

// Count returns the number of tokens in the given text.
// This is useful for checking token limits before processing large texts.
//
// Returns:
//   - (count, nil) where count is the number of tokens
//   - (0, ErrInvalidTokenizer) if the tokenizer is invalid
func (t *Tokenizer) Count(text string) (int, error) {
	if t.ptr == nil {
		return 0, ErrInvalidTokenizer
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	count := C.bpe_count(t.ptr, cText)
	return int(count), nil
}

// CountTillLimit returns the token count if it's below the given limit,
// otherwise returns -1. This is more efficient than Count when you only
// need to know if text exceeds a token limit.
//
// Parameters:
//   - text: The input text to count tokens for
//   - limit: Maximum number of tokens to count up to
//
// Returns:
//   - (count, nil) where count is the number of tokens if below limit
//   - (-1, nil) if the token count exceeds the limit
//   - (0, ErrInvalidTokenizer) if the tokenizer is invalid
func (t *Tokenizer) CountTillLimit(text string, limit int) (int, error) {
	if t.ptr == nil {
		return 0, ErrInvalidTokenizer
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	count := C.bpe_count_till_limit(t.ptr, cText, C.size_t(limit))
	if count == ^C.size_t(0) {  
		return -1, nil
	}
	return int(count), nil
}

// Encode converts the text into a sequence of token IDs.
// The function handles Unicode characters correctly and is thread-safe.
//
// Returns:
//   - (tokens, nil) where tokens is a slice of uint32 token IDs
//   - (nil, ErrInvalidTokenizer) if the tokenizer is invalid
//   - (nil, ErrEncoding) if encoding fails
//   - (empty slice, nil) for empty input text
func (t *Tokenizer) Encode(text string) ([]uint32, error) {
	if t.ptr == nil {
		return nil, ErrInvalidTokenizer
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	var tokenCount C.size_t
	tokensPtr := C.bpe_encode(t.ptr, cText, &tokenCount)
	if tokensPtr == nil && tokenCount > 0 {
		return nil, ErrEncoding
	}

	if tokenCount == 0 {
		return []uint32{}, nil
	}

	// Create a slice from the C array without copying
	tokens := unsafe.Slice((*uint32)(unsafe.Pointer(tokensPtr)), int(tokenCount))
	
	// Create a new Go slice and copy the data
	result := make([]uint32, int(tokenCount))
	copy(result, tokens)

	// Free the original C array
	C.free(unsafe.Pointer(tokensPtr))

	return result, nil
}

// Decode converts a sequence of token IDs back into text.
// The function correctly handles Unicode characters and special tokens.
//
// Returns:
//   - (text, nil) where text is the decoded string
//   - ("", ErrInvalidTokenizer) if the tokenizer is invalid
//   - ("", ErrDecoding) if decoding fails
//   - ("", nil) for empty input tokens
func (t *Tokenizer) Decode(tokens []uint32) (string, error) {
	if t.ptr == nil {
		return "", ErrInvalidTokenizer
	}
	if len(tokens) == 0 {
		return "", nil
	}

	cTokens := (*C.uint32_t)(unsafe.Pointer(&tokens[0]))
	cText := C.bpe_decode(t.ptr, cTokens, C.size_t(len(tokens)))
	if cText == nil {
		return "", ErrDecoding
	}
	defer C.free(unsafe.Pointer(cText))

	return C.GoString(cText), nil
}

// free releases the underlying Rust tokenizer
func (t *Tokenizer) free() {
	if t.ptr != nil {
		C.bpe_free(t.ptr)
		t.ptr = nil
	}
}
