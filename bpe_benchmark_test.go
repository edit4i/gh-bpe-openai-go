package bpe

import (
	"testing"

	"github.com/pkoukk/tiktoken-go"
)

var benchmarkTexts = []struct {
	name string
	text string
}{
	{"Short", "Hello, world!"},
	{"Medium", "This is a medium length text that contains multiple sentences. It should give us a good idea of tokenization performance."},
	{"Long", `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
		Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 
		Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.`},
	{"Unicode", "Hello 👋 World 🌍! This test includes emojis 🎉 and other Unicode characters: 你好, こんにちは, Привет"},
}

func BenchmarkCL100kTokenizer(b *testing.B) {
	tok, err := NewCL100kTokenizer()
	if err != nil {
		b.Fatalf("Failed to create CL100k tokenizer: %v", err)
	}

	for _, tc := range benchmarkTexts {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := tok.Encode(tc.text)
				if err != nil {
					b.Fatalf("Encode failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkO200kTokenizer(b *testing.B) {
	tok, err := NewO200kTokenizer()
	if err != nil {
		b.Fatalf("Failed to create O200k tokenizer: %v", err)
	}

	for _, tc := range benchmarkTexts {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := tok.Encode(tc.text)
				if err != nil {
					b.Fatalf("Encode failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkTiktokenCL100k(b *testing.B) {
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		b.Fatalf("Failed to create CL100k tokenizer: %v", err)
	}

	for _, tc := range benchmarkTexts {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = enc.Encode(tc.text, nil, nil)
			}
		})
	}
}

func BenchmarkTiktokenO200k(b *testing.B) {
	enc, err := tiktoken.GetEncoding("o200k_base")
	if err != nil {
		b.Fatalf("Failed to create O200k tokenizer: %v", err)
	}

	for _, tc := range benchmarkTexts {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = enc.Encode(tc.text, nil, nil)
			}
		})
	}
}
