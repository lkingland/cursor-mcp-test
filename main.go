package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	server "cursor-mcp-test/pkg"
)

const Logfile = "/tmp/cursor-mcp-test.log"

func main() {
	logFile, err := os.Create(Logfile)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log.Printf("MCP communication will be logged to %s", Logfile)

	transport := &mcp.LoggingTransport{
		Transport: &mcp.StdioTransport{},
		Writer:    logFile,
	}

	if err := server.New().Run(context.Background(), transport); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
