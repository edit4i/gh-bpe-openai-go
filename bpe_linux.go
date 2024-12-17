//go:build linux && amd64
// +build linux,amd64

package bpe

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -lbpe_openai_ffi -Wl,-rpath=${SRCDIR}/lib/linux_amd64
import "C"
