package main

import (
	"context"
	"log"
	"os"

	"github.com/Warashi/go-modelcontextprotocol/mcp"
)

const (
	PIPECD_HOST_ENV     = "PIPECD_HOST"
	PIPECD_API_KEY_ENV  = "PIPECD_API_KEY"
	PIPECD_INSECURE_ENV = "PIPECD_INSECURE"
)

func main() {
	addr := os.Getenv(PIPECD_HOST_ENV)
	key := os.Getenv(PIPECD_API_KEY_ENV)
	insecure := os.Getenv(PIPECD_INSECURE_ENV) == "true"

	s, err := NewServer(context.Background(), addr, key, insecure)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	opts := []mcp.ServerOption{
		mcp.WithResourceReader(s),
		mcp.WithResourceTemplate(resourceTemplateApplications),
		mcp.WithResourceTemplate(resourceTemplateDeployments),
		mcp.WithTool(s.listApplicationsTool()),
	}

	server := mcp.NewStdioServer("PipeCD MCP Server", "0.0.1", opts...)

	if err := server.Serve(context.Background()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
