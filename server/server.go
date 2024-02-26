package server

import (
	"errors"
	"fmt"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type Server struct {
	listener net.Listener
	srv      *grpc.Server
}

func NewServer() *Server {
	s := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(insecure.NewCredentials()),
	)
	notesv1.RegisterNotesServiceServer(s, NewNotesService())

	srv := &Server{
		srv: s,
	}

	return srv
}

func (s *Server) Listen(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %e", err)
	}
	s.listener = lis
	return s.srv.Serve(s.listener)
}

func (s *Server) Close() error {
	s.srv.GracefulStop()
	if err := s.listener.Close(); !errors.Is(err, net.ErrClosed) {
		return err
	}
	return nil
}
