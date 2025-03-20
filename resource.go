package main

import (
	"context"
	"net/url"
	"path"

	"github.com/Warashi/go-modelcontextprotocol/jsonrpc2"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"google.golang.org/protobuf/encoding/protojson"
)

var resourceTemplateApplications = mcp.ResourceTemplate{
	URITemplate: "pipecd://applications/{applicationId}",
	Name:        "Application",
	Description: "An application managed by PipeCD",
	MimeType:    "application/json",
}

var resourceTemplateDeployments = mcp.ResourceTemplate{
	URITemplate: "pipecd://deployments/{deploymentId}",
	Name:        "Deployment",
	Description: "A deployment managed by PipeCD",
	MimeType:    "application/json",
}

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
	}

	return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "unsupported resource type", struct{}{})
}

func (s *Server) readApplication(ctx context.Context, u *url.URL) (*mcp.Result[mcp.ReadResourceResultData], error) {
	id := path.Base(u.Path)
	if id == "" {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "missing application ID", struct{}{})
	}

	app, err := s.client.GetApplication(ctx, &apiservice.GetApplicationRequest{
		ApplicationId: id,
	})
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInternalError, "failed to get application", struct{}{})
	}

	b, err := protojson.Marshal(app)
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInternalError, "failed to marshal application", struct{}{})
	}

	return &mcp.Result[mcp.ReadResourceResultData]{
		Data: mcp.ReadResourceResultData{
			Contents: []mcp.IsResourceContents{
				mcp.TextResourceContents{
					URI:      u.String(),
					MimeType: "application/json",
					Text:     string(b),
				},
			},
		},
	}, nil
}

func (s *Server) readDeployment(ctx context.Context, u *url.URL) (*mcp.Result[mcp.ReadResourceResultData], error) {
	id := path.Base(u.Path)
	if id == "" {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "missing deployment ID", struct{}{})
	}

	deployment, err := s.client.GetDeployment(ctx, &apiservice.GetDeploymentRequest{
		DeploymentId: id,
	})
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInternalError, "failed to get deployment", struct{}{})
	}

	b, err := protojson.Marshal(deployment)
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInternalError, "failed to marshal deployment", struct{}{})
	}

	return &mcp.Result[mcp.ReadResourceResultData]{
		Data: mcp.ReadResourceResultData{
			Contents: []mcp.IsResourceContents{
				mcp.TextResourceContents{
					URI:      u.String(),
					MimeType: "application/json",
					Text:     string(b),
				},
			},
		},
	}, nil
}
