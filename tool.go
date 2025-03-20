package main

import (
	"context"
	"encoding/json"

	"github.com/Warashi/go-modelcontextprotocol/jsonschema"
	"github.com/Warashi/go-modelcontextprotocol/mcp"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
)

type listApplicationsResponse struct {
	Applications []*model.Application
	NextCursor   string
}

func (r *listApplicationsResponse) MarshalJSON() ([]byte, error) {
	apps := make([]json.RawMessage, 0, len(r.Applications))
	for _, app := range r.Applications {
		b, err := protojson.Marshal(app)
		if err != nil {
			return nil, err
		}
		apps = append(apps, b)
	}

	return json.Marshal(struct {
		Applications []json.RawMessage `json:"applications"`
		NextCursor   string            `json:"nextCursor"`
	}{
		Applications: apps,
		NextCursor:   r.NextCursor,
	})
}

type listApplicationsRequest struct {
	Cursor string `json:"cursor"`
}

func (s *Server) listApplicationsTool() mcp.Tool[*listApplicationsRequest] {
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

func (s *Server) listApplications(ctx context.Context, request *listApplicationsRequest) ([]any, error) {
	response, err := s.client.ListApplications(ctx, &apiservice.ListApplicationsRequest{
		Cursor: request.Cursor,
	})
	if err != nil {
		return nil, err
	}
	return []any{&listApplicationsResponse{
		Applications: response.GetApplications(),
		NextCursor:   response.GetCursor(),
	}}, nil

}
