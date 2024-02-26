package server

import (
	"errors"
	"fmt"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server struct {
	listener net.Listener
	grpcSrv  *grpc.Server
}

func NewServer(c credentials.TransportCredentials) *Server {
	s := grpc.NewServer(
		grpc.Creds(c),
	)
	srv := &Server{
		grpcSrv: s,
	}
	srv.registerServices()

	return srv
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
	if err := s.listener.Close(); !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("error closing server: %v", err)
	}
	return nil
}

func (s *Server) registerServices() {
	notesv1.RegisterNotesServiceServer(s.grpcSrv, NewNotesService())
}
