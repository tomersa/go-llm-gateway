package main

import (
	"log"

	"github.com/tomersa/llm-gateway/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
