//go:build windows && amd64
// +build windows,amd64

package bpe

// #cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -lbpe_openai_ffi
// #cgo LDFLAGS: -Wl,--enable-auto-import
import "C"
