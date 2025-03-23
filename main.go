package main

import (
	"context"
	"log"
	"os"

	"github.com/Warashi/go-modelcontextprotocol/mcp"
)

const (
	PIPECD_HOST_ENV         = "PIPECD_HOST"
	PIPECD_API_KEY_ENV      = "PIPECD_API_KEY"
	PIPECD_API_KEY_FILE_ENV = "PIPECD_API_KEY_FILE"
	PIPECD_INSECURE_ENV     = "PIPECD_INSECURE"
)

func main() {
	addr := os.Getenv(PIPECD_HOST_ENV)
	key := os.Getenv(PIPECD_API_KEY_ENV)
	if key == "" {
		keyFile := os.Getenv(PIPECD_API_KEY_FILE_ENV)
		if keyFile != "" {
			b, err := os.ReadFile(keyFile)
			if err != nil {
				log.Fatalf("failed to read api key file: %v", err)
			}
			key = string(b)
		}
	}

	insecure := os.Getenv(PIPECD_INSECURE_ENV) == "true"

	if addr == "" {
		log.Fatalln("PIPECD_HOST is required")
	}
	if key == "" {
		log.Fatalln("PIPECD_API_KEY or PIPECD_API_KEY_FILE is required")
	}

	s, err := NewServer(context.Background(), addr, key, insecure)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	mux := mcp.NewResourceReaderMux()
	mux.HandleFunc("pipecd://applications/{applicationId}", s.readApplication)
	mux.HandleFunc("pipecd://deployments/{deploymentId}", s.readDeployment)
	mux.HandleFunc("pipecd://deployments/{deploymentId}/logs/{stageId}", s.readDeploymentStageLogs)

	opts := []mcp.ServerOption{
		mcp.WithResourceReader(mux),
		mcp.WithResourceTemplate(resourceTemplateApplications),
		mcp.WithResourceTemplate(resourceTemplateDeployments),
		mcp.WithResourceTemplate(resourceTemplateDeploymentStageLogs),
		mcp.WithTool(s.listApplicationsTool()),
		mcp.WithTool(s.listDeploymentsTool()),
		mcp.WithTool(s.getDeploymentTool()),
		mcp.WithTool(s.getDeploymentStageLogsTool()),
	}

	server := mcp.NewStdioServer("PipeCD MCP Server", "0.0.1", opts...)

	if err := server.Serve(context.Background()); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
