package server

import (
	"errors"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"time"
)

type Server struct {
	listener  net.Listener
	grpcSrv   *grpc.Server
	healthSrv *health.Server
}

func NewServer(c credentials.TransportCredentials) (*Server, error) {
	s := grpc.NewServer(
		grpc.Creds(c),
	)
	srv := &Server{
		grpcSrv:   s,
		healthSrv: health.NewServer(),
	}

	v, err := protovalidate.New()
	if err != nil {
		return srv, fmt.Errorf("failed to initialize validator: %v", err)
	}

	grpc_health_v1.RegisterHealthServer(srv.grpcSrv, srv.healthSrv)
	srv.registerServices(v)

	return srv, nil
}

func (s *Server) Listen(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s.listener = lis
	return s.grpcSrv.Serve(s.listener)
}

func (s *Server) Close() error {
	s.healthSrv.Shutdown()
	time.Sleep(5 * time.Second) // Set a timeout in which all services will be marked as unhealthy

	s.grpcSrv.GracefulStop()

	err := s.listener.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("error closing server: %v", err)
	}

	return nil
}

func (s *Server) registerServices(v *protovalidate.Validator) {
	notesv1.RegisterNotesServiceServer(s.grpcSrv, NewNotesService(v))
	s.healthSrv.SetServingStatus("notes.v1.NotesService", grpc_health_v1.HealthCheckResponse_SERVING)
}
