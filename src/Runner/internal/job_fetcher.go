package internal

import (
	"context"

	"github.com/Jakkoble/HexaTask/src/Runner/pb"
)

type IJobFetcher interface {
	FetchCommands(ctx context.Context, jobID string) ([]string, error)
}

type GrpcJobFetcher struct {
	client pb.RunnerServiceClient
}

func NewGrpcJobFetcher(client pb.RunnerServiceClient) *GrpcJobFetcher {
	return &GrpcJobFetcher{client: client}
}

func (f *GrpcJobFetcher) FetchCommands(ctx context.Context, jobID string) ([]string, error) {
	res, err := f.client.GetJobDetails(ctx, &pb.GetJobDetailsRequest{JobId: jobID})
	if err != nil {
		return nil, err
	}

	return res.Commands, nil
}
