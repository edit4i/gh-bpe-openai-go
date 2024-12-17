package main

import (
	"fmt"
	"log"

	bpe "github.com/edit4i/gh-bpe-openai-go"
)

func main() {
	// Initialize the tokenizer with CL100k model
	tokenizer, err := bpe.NewCL100kTokenizer()
	if err != nil {
		log.Fatal(err)
	}

	// Example texts to tokenize
	texts := []string{
		"Hello, this is a test of the OpenAI tokenizer!",
		"Hello 👋 World 🌍", // Unicode example
		"",                  // Empty string
	}

	for _, text := range texts {
		fmt.Printf("\nProcessing text: %q\n", text)

		// Count tokens
		count, err := tokenizer.Count(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Token count: %d\n", count)

		// Encode the text
		tokens, err := tokenizer.Encode(text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Tokens: %v\n", tokens)

		// Decode back to text
		decoded, err := tokenizer.Decode(tokens)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Decoded: %q\n", decoded)

		// Verify round-trip
		if decoded != text {
			fmt.Printf("Warning: Round-trip encoding/decoding produced different result!\n")
		}
	}
}
