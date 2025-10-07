package main

import (
	"fmt"
	"io"
	"os"

	"github.com/james/tasks-planner/internal/canonjson"
	"github.com/james/tasks-planner/internal/hash"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Missing required argument <file-path>")
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/tasksd <file-path>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open file '%s': %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	inputBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	canonicalBytes, err := canonjson.ToCanonicalJSON(inputBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to process JSON: %v\n", err)
		os.Exit(1)
	}

	hashString := hash.HashCanonicalBytes(canonicalBytes)

	fmt.Println("--- Canonical JSON ---")
	fmt.Print(string(canonicalBytes))
	fmt.Println("--- SHA-256 Hash ---")
	fmt.Println(hashString)
}
