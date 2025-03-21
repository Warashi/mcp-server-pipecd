package main

import (
	"context"
	"net/url"
	"path"
	"strings"

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

var resourceTemplateDeploymentStageLogs = mcp.ResourceTemplate{
	URITemplate: "pipecd://deployment-stage-logs/{deploymentId}/{stageId}",
	Name:        "Deployment Stage Logs",
	Description: "Logs of a deployment stage managed by PipeCD",
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
	case "deployment-stage-logs":
		return s.readDeploymentStageLogs(ctx, u)
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

func (s *Server) readDeploymentStageLogs(ctx context.Context, u *url.URL) (*mcp.Result[mcp.ReadResourceResultData], error) {
	fields := strings.Split(strings.TrimPrefix(strings.TrimSuffix(u.Path, "/"), "/"), "/")
	if len(fields) != 2 {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInvalidParams, "invalid deployment stage logs URI", struct{}{})
	}

	id := fields[0]
	stageID := fields[1]

	logs, err := s.client.ListStageLogs(ctx, &apiservice.ListStageLogsRequest{
		DeploymentId: id,
	})
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInternalError, "failed to get deployment stage logs", struct{}{})
	}

	log := logs.GetStageLogs()[stageID]
	if log == nil {
		return &mcp.Result[mcp.ReadResourceResultData]{
			Data: mcp.ReadResourceResultData{
				Contents: []mcp.IsResourceContents{
					mcp.TextResourceContents{
						URI:      u.String(),
						MimeType: "application/json",
						Text:     "{ \"error\": \"no logs found\" }",
					},
				},
			},
		}, nil
	}

	b, err := protojson.Marshal(log)
	if err != nil {
		return nil, jsonrpc2.NewError(jsonrpc2.CodeInternalError, "failed to marshal deployment stage logs", struct{}{})
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
