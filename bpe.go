package bpe

// #cgo LDFLAGS: -L${SRCDIR}/build -lbpe_openai_ffi
// #cgo CFLAGS: -I${SRCDIR}/include
// #include <bpe_openai.h>
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// Tokenizer represents a BPE tokenizer instance
type Tokenizer struct {
	ptr *C.TokenizerHandle
}

// Error definitions
var (
	ErrInvalidTokenizer = errors.New("invalid tokenizer")
	ErrEncoding         = errors.New("encoding error")
	ErrDecoding         = errors.New("decoding error")
)

// NewCL100kTokenizer creates a new CL100k tokenizer instance
func NewCL100kTokenizer() (*Tokenizer, error) {
	ptr := C.bpe_cl100k_base()
	if ptr == nil {
		return nil, ErrInvalidTokenizer
	}
	t := &Tokenizer{ptr: ptr}
	runtime.SetFinalizer(t, (*Tokenizer).free)
	return t, nil
}

// NewO200kTokenizer creates a new O200k tokenizer instance
func NewO200kTokenizer() (*Tokenizer, error) {
	ptr := C.bpe_o200k_base()
	if ptr == nil {
		return nil, ErrInvalidTokenizer
	}
	t := &Tokenizer{ptr: ptr}
	runtime.SetFinalizer(t, (*Tokenizer).free)
	return t, nil
}

// Count returns the number of tokens in the given text
func (t *Tokenizer) Count(text string) (int, error) {
	if t.ptr == nil {
		return 0, ErrInvalidTokenizer
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	count := C.bpe_count(t.ptr, cText)
	return int(count), nil
}

// CountTillLimit returns the token count if it's below the given limit, otherwise returns -1
func (t *Tokenizer) CountTillLimit(text string, limit int) (int, error) {
	if t.ptr == nil {
		return 0, ErrInvalidTokenizer
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	count := C.bpe_count_till_limit(t.ptr, cText, C.size_t(limit))
	if count == C.size_t(-1) {
		return -1, nil
	}
	return int(count), nil
}

// Encode converts the text into a sequence of token IDs
func (t *Tokenizer) Encode(text string) ([]uint32, error) {
	if t.ptr == nil {
		return nil, ErrInvalidTokenizer
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	var tokenCount C.size_t
	tokensPtr := C.bpe_encode(t.ptr, cText, &tokenCount)
	if tokensPtr == nil {
		return nil, ErrEncoding
	}
	defer C.free(unsafe.Pointer(tokensPtr))

	tokens := make([]uint32, int(tokenCount))
	src := unsafe.Slice((*uint32)(unsafe.Pointer(tokensPtr)), int(tokenCount))
	copy(tokens, src)

	return tokens, nil
}

// Decode converts a sequence of token IDs back into text
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
