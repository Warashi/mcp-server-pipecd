package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Warashi/go-modelcontextprotocol/jsonschema"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
)

type listDeploymentsRequest struct {
	ApplicationID string `json:"applicationId"`
	Cursor        string `json:"cursor"`
}

type listDeploymentsResponse struct {
	Deployments []*deployment `json:"deployments"`
	NextCursor  string        `json:"nextCursor"`
}

type getDeploymentRequest struct {
	ID string `json:"id"`
}

type getDeploymentStageLogsRequest struct {
	DeploymentID string `json:"deploymentId"`
	StageID      string `json:"stageId"`
}

type deployment struct {
	ID            string               `json:"id"`
	ApplicationID string               `json:"applicationId"`
	Labels        map[string]string    `json:"labels"`
	Summary       string               `json:"summary"`
	Status        string               `json:"status"`
	StatusReason  string               `json:"statusReason"`
	Artifacts     []deploymentArtifact `json:"artifacts"`
	Stages        []stage              `json:"stages"`
}

type deploymentArtifact struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

type stage struct {
	ID           string `json:"id"`
	StageName    string `json:"stageName"`
	Status       string `json:"status"`
	StatusReason string `json:"statusReason"`
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

	stages := make([]stage, 0, len(d.GetStages()))
	for _, s := range d.GetStages() {
		stages = append(stages, stage{
			ID:           s.GetId(),
			StageName:    s.GetName(),
			Status:       s.GetStatus().String(),
			StatusReason: s.GetStatusReason(),
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
		Stages:        stages,
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
	var applicationIDs []string
	if request.ApplicationID != "" {
		applicationIDs = []string{request.ApplicationID}
	}

	response, err := s.client.ListDeployments(ctx, &apiservice.ListDeploymentsRequest{
		ApplicationIds: applicationIDs,
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

func (s *Server) getDeploymentTool() mcp.Tool[*getDeploymentRequest, *mcp.EmbeddedResource] {
	return mcp.NewToolFunc(
		"GetDeployment",
		"Get a deployment managed by PipeCD",
		jsonschema.Object{
			Properties: map[string]jsonschema.Schema{
				"id": jsonschema.String{},
			},
			Required: []string{"id"},
		},
		s.getDeployment,
	)
}

func (s *Server) getDeployment(ctx context.Context, request *getDeploymentRequest) (*mcp.EmbeddedResource, error) {
	response, err := s.client.GetDeployment(ctx, &apiservice.GetDeploymentRequest{
		DeploymentId: request.ID,
	})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(newDeployment(response.Deployment))
	if err != nil {
		return nil, err
	}

	return &mcp.EmbeddedResource{
		Resource: mcp.TextResourceContents{
			URI:      fmt.Sprintf("pipecd://deployment/%s", response.Deployment.GetId()),
			MimeType: "application/json",
			Text:     string(b),
		},
	}, nil
}

func (s *Server) getDeploymentStageLogsTool() mcp.Tool[*getDeploymentStageLogsRequest, *mcp.EmbeddedResource] {
	return mcp.NewToolFunc(
		"GetDeploymentStageLogs",
		"Get deployment stage logs managed by PipeCD",
		jsonschema.Object{
			Properties: map[string]jsonschema.Schema{
				"deploymentId": jsonschema.String{},
				"stageId":      jsonschema.String{},
			},
			Required: []string{"deploymentId", "stageId"},
		},
		s.getDeploymentStageLogs,
	)
}

func (s *Server) getDeploymentStageLogs(ctx context.Context, request *getDeploymentStageLogsRequest) (*mcp.EmbeddedResource, error) {
	response, err := s.client.ListStageLogs(ctx, &apiservice.ListStageLogsRequest{
		DeploymentId: request.DeploymentID,
	})
	if err != nil {
		return nil, err
	}

	log := response.GetStageLogs()[request.StageID]
	if log == nil {
		return nil, fmt.Errorf("no logs found")
	}

	b, err := protojson.Marshal(log)
	if err != nil {
		return nil, err
	}

	return &mcp.EmbeddedResource{
		Resource: mcp.TextResourceContents{
			URI:      fmt.Sprintf("pipecd://deployment-stage-logs/%s/%s", request.DeploymentID, request.StageID),
			MimeType: "application/json",
			Text:     string(b),
		},
	}, nil
}
