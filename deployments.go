package main

import (
	"context"

	"github.com/Warashi/go-modelcontextprotocol/jsonschema"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type listDeploymentsRequest struct {
	ApplicationID string `json:"applicationId"`
	Cursor        string `json:"cursor"`
}

type listDeploymentsResponse struct {
	Deployments []*deployment `json:"deployments"`
	NextCursor  string        `json:"nextCursor"`
}

type deployment struct {
	ID            string               `json:"id"`
	ApplicationID string               `json:"applicationId"`
	Labels        map[string]string    `json:"labels"`
	Summary       string               `json:"summary"`
	Status        string               `json:"status"`
	StatusReason  string               `json:"statusReason"`
	Artifacts     []deploymentArtifact `json:"artifacts"`
}

type deploymentArtifact struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

func newDeployment(d *model.Deployment) *deployment {
	artifacts := make([]deploymentArtifact, 0, len(d.GetVersions()))
	for _, a := range d.GetVersions() {
		artifacts = append(artifacts, deploymentArtifact{
			Kind:    a.GetKind().String(),
			Name:    a.GetName(),
			Version: a.GetVersion(),
			URL:     a.GetUrl(),
		})
	}

	return &deployment{
		ID:            d.GetId(),
		ApplicationID: d.GetApplicationId(),
		Labels:        d.GetLabels(),
		Summary:       d.GetSummary(),
		Status:        d.GetStatus().String(),
		StatusReason:  d.GetStatusReason(),
		Artifacts:     artifacts,
	}
}

func (s *Server) listDeploymentsTool() mcp.Tool[*listDeploymentsRequest, *listDeploymentsResponse] {
	return mcp.NewToolFunc(
		"ListDeployments",
		"List deployments managed by PipeCD",
		jsonschema.Object{
			Properties: map[string]jsonschema.Schema{
				"cursor":        jsonschema.String{},
				"applicationId": jsonschema.String{},
			},
		},
		s.listDeployments,
	)
}

func (s *Server) listDeployments(ctx context.Context, request *listDeploymentsRequest) (*listDeploymentsResponse, error) {
	response, err := s.client.ListDeployments(ctx, &apiservice.ListDeploymentsRequest{
		ApplicationIds: []string{request.ApplicationID},
		Cursor:         request.Cursor,
	})
	if err != nil {
		return nil, err
	}

	deployments := make([]*deployment, 0, len(response.GetDeployments()))
	for _, d := range response.GetDeployments() {
		deployments = append(deployments, newDeployment(d))
	}

	return &listDeploymentsResponse{
		Deployments: deployments,
		NextCursor:  response.GetCursor(),
	}, nil
}
