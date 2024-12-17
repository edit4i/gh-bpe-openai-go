//go:build darwin && amd64
// +build darwin,amd64

package bpe

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_amd64 -lbpe_openai_ffi -Wl,-rpath,${SRCDIR}/lib/darwin_amd64
import "C"
