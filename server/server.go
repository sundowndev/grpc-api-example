package server

import (
	"errors"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	healthv1 "github.com/sundowndev/grpc-api-example/proto/health/v1"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server struct {
	listener net.Listener
	grpcSrv  *grpc.Server
}

func NewServer(c credentials.TransportCredentials) (*Server, error) {
	s := grpc.NewServer(
		grpc.Creds(c),
	)
	srv := &Server{
		grpcSrv: s,
	}

	v, err := protovalidate.New()
	if err != nil {
		return srv, fmt.Errorf("failed to initialize validator: %v", err)
	}

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
	s.grpcSrv.GracefulStop()

	err := s.listener.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("error closing server: %v", err)
	}

	return nil
}

func (s *Server) registerServices(v *protovalidate.Validator) {
	healthv1.RegisterHealthServiceServer(s.grpcSrv, NewHealthService())
	notesv1.RegisterNotesServiceServer(s.grpcSrv, NewNotesService(v))
}
