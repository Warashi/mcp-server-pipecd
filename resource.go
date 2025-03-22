package main

import (
	"context"
	"net/url"

	"github.com/Warashi/go-modelcontextprotocol/jsonrpc2"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
)

func (s *Server) ReadResource(ctx context.Context, request *mcp.Request[mcp.ReadResourceRequestParams]) (*mcp.Result[mcp.ReadResourceResultData], error) {
	u, err := url.Parse(request.Params.URI)
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "failed to parse URI", struct{}{})
	}
	if u.Scheme != "pipecd" {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "unsupported URI scheme", struct{}{})
	}

	switch u.Host {
	case "applications":
		return s.readApplication(ctx, u)
	case "deployments":
		return s.readDeployment(ctx, u)
	case "deployment-stage-logs":
		return s.readDeploymentStageLogs(ctx, u)
	}

	return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "unsupported resource type", struct{}{})
}
