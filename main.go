package main

import (
	"context"
	"log"

	"github.com/Warashi/go-modelcontextprotocol/mcp"
)

func main() {
	opts := []mcp.ServerOption{}

	server := mcp.NewStdioServer("PipeCD MCP Server", "0.0.1", opts...)

	if err := server.Serve(context.Background()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
