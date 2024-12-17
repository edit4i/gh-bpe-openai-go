package bpe

import (
	"testing"
)

func TestCL100kTokenizer(t *testing.T) {
	tok, err := NewCL100kTokenizer()
	if err != nil {
		t.Fatalf("Failed to create CL100k tokenizer: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantErr  bool
		wantText string // for encode/decode roundtrip
	}{
		{
			name:     "Simple text",
			input:    "Hello, world!",
			wantLen:  4,
			wantErr:  false,
			wantText: "Hello, world!",
		},
		{
			name:     "Empty string",
			input:    "",
			wantLen:  0,
			wantErr:  false,
			wantText: "",
		},
		{
			name:     "Unicode text",
			input:    "Hello 👋 World 🌍",
			wantLen:  7,
			wantErr:  false,
			wantText: "Hello 👋 World 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Count
			count, err := tok.Count(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if count != tt.wantLen {
				t.Errorf("Count() = %v, want %v", count, tt.wantLen)
			}

			// Test CountTillLimit
			count, err = tok.CountTillLimit(tt.input, tt.wantLen+1)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountTillLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if count != tt.wantLen {
				t.Errorf("CountTillLimit() = %v, want %v", count, tt.wantLen)
			}

			// Test Encode/Decode roundtrip
			tokens, err := tok.Encode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decoded, err := tok.Decode(tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if decoded != tt.wantText {
				t.Errorf("Encode/Decode roundtrip failed: got %q, want %q", decoded, tt.wantText)
			}
		})
	}
}

func TestO200kTokenizer(t *testing.T) {
	tok, err := NewO200kTokenizer()
	if err != nil {
		t.Fatalf("Failed to create O200k tokenizer: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantErr  bool
		wantText string // for encode/decode roundtrip
	}{
		{
			name:     "Simple text",
			input:    "Hello, world!",
			wantLen:  4,
			wantErr:  false,
			wantText: "Hello, world!",
		},
		{
			name:     "Empty string",
			input:    "",
			wantLen:  0,
			wantErr:  false,
			wantText: "",
		},
		{
			name:     "Unicode text",
			input:    "Hello 👋 World 🌍",
			wantLen:  6,
			wantErr:  false,
			wantText: "Hello 👋 World 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Count
			count, err := tok.Count(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if count != tt.wantLen {
				t.Errorf("Count() = %v, want %v", count, tt.wantLen)
			}

			// Test Encode/Decode roundtrip
			tokens, err := tok.Encode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decoded, err := tok.Decode(tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if decoded != tt.wantText {
				t.Errorf("Encode/Decode roundtrip failed: got %q, want %q", decoded, tt.wantText)
			}
		})
	}
}
