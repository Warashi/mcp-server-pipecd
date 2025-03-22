package main

import (
	"context"
	"net/url"
	"path"

	"github.com/Warashi/go-modelcontextprotocol/jsonrpc2"
	"github.com/Warashi/go-modelcontextprotocol/jsonschema"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
)

var resourceTemplateApplications = mcp.ResourceTemplate{
	URITemplate: "pipecd://applications/{applicationId}",
	Name:        "Application",
	Description: "An application managed by PipeCD",
	MimeType:    "application/json",
}

type listApplicationsRequest struct {
	Cursor string `json:"cursor"`
}

type listApplicationsResponse struct {
	Applications []*application `json:"applications"`
	NextCursor   string         `json:"nextCursor"`
}

type application struct {
	ID                         string                `json:"id"`
	Name                       string                `json:"name"`
	PipedID                    string                `json:"pipedId"`
	ProjectID                  string                `json:"projectId"`
	DeployTargets              []string              `json:"deployTargets"`
	Description                string                `json:"description"`
	Labels                     map[string]string     `json:"labels"`
	LastSuccessfulDeploymentID string                `json:"lastSuccessfulDeploymentID"`
	SyncState                  *applicationSyncState `json:"syncState"`
}

type applicationSyncState struct {
	Status           string `json:"status"`
	ShortReason      string `json:"shortReason"`
	Reason           string `json:"reason"`
	HeadDeploymentID string `json:"headDeploymentId"`
}

func newApplication(a *model.Application) *application {
	return &application{
		ID:                         a.GetId(),
		Name:                       a.GetName(),
		PipedID:                    a.GetPipedId(),
		ProjectID:                  a.GetProjectId(),
		DeployTargets:              a.GetDeployTargets(),
		Description:                a.GetDescription(),
		Labels:                     a.GetLabels(),
		LastSuccessfulDeploymentID: a.GetMostRecentlySuccessfulDeployment().GetDeploymentId(),
		SyncState: &applicationSyncState{
			Status:           a.GetSyncState().GetStatus().String(),
			ShortReason:      a.GetSyncState().GetShortReason(),
			Reason:           a.GetSyncState().GetReason(),
			HeadDeploymentID: a.GetSyncState().GetHeadDeploymentId(),
		},
	}
}

func (s *Server) listApplicationsTool() mcp.Tool[*listApplicationsRequest, *listApplicationsResponse] {
	return mcp.NewToolFunc(
		"ListApplications",
		"List applications managed by PipeCD",
		jsonschema.Object{
			Properties: map[string]jsonschema.Schema{
				"cursor": jsonschema.String{},
			},
		},
		s.listApplications,
	)
}

func (s *Server) listApplications(ctx context.Context, request *listApplicationsRequest) (*listApplicationsResponse, error) {
	response, err := s.client.ListApplications(ctx, &apiservice.ListApplicationsRequest{
		Cursor: request.Cursor,
	})
	if err != nil {
		return nil, err
	}

	apps := make([]*application, 0, len(response.Applications))
	for _, a := range response.Applications {
		apps = append(apps, newApplication(a))
	}

	return &listApplicationsResponse{
		Applications: apps,
		NextCursor:   response.GetCursor(),
	}, nil
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
