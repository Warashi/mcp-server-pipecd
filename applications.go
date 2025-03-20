package main

import (
	"context"

	"github.com/Warashi/go-modelcontextprotocol/jsonschema"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

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
