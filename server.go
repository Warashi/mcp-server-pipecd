package main

import (
	"context"
	"crypto/tls"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	client apiservice.Client
}

func NewServer(ctx context.Context, addr, key string, insecure bool) (*Server, error) {
	creds := rpcclient.NewPerRPCCredentials(key, rpcauth.APIKeyCredentials, !insecure)

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithPerRPCCredentials(creds),
	}

	if !insecure {
		options = append(options, rpcclient.WithTransportCredentials(credentials.NewTLS(new(tls.Config))))
	}

	client, err := apiservice.NewClient(ctx, addr, options...)
	if err != nil {
		return nil, err
	}

	return &Server{
		client: client,
	}, nil
}
