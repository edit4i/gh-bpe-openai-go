//go:build darwin && arm64
// +build darwin,arm64

package bpe

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lbpe_openai_ffi -Wl,-rpath,${SRCDIR}/lib/darwin_arm64
import "C"
