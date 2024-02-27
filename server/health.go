package server

import (
	"context"
	healthv1 "github.com/sundowndev/grpc-api-example/proto/health/v1"
)

type HealthService struct {
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) GetHealth(_ context.Context, _ *healthv1.GetHealthRequest) (*healthv1.GetHealthResponse, error) {
	return &healthv1.GetHealthResponse{Ok: true}, nil
}
