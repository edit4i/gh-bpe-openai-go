//go:build linux && arm64
// +build linux,arm64

package bpe

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_arm64 -lbpe_openai_ffi -Wl,-rpath=${SRCDIR}/lib/linux_arm64
import "C"
