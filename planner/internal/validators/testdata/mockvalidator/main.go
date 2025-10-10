package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type payload struct {
	Tasks any `json:"tasks"`
}

type result struct {
	Status    string `json:"status"`
	Validator string `json:"validator"`
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read stdin: %v", err)
	}
	var p payload
	if err := json.Unmarshal(data, &p); err != nil {
		log.Fatalf("unmarshal payload: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(result{Status: "pass", Validator: "mock"}); err != nil {
		log.Fatalf("encode result: %v", err)
	}
}
