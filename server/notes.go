package server

import (
	"context"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"github.com/gofrs/uuid"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"sync"
)

type NotesService struct {
	notesv1.UnimplementedNotesServiceServer
	mu        *sync.RWMutex
	notes     []*notesv1.Note
	validator *protovalidate.Validator
}

func NewNotesService(v *protovalidate.Validator) *NotesService {
	return &NotesService{
		mu:        &sync.RWMutex{},
		notes:     make([]*notesv1.Note, 0),
		validator: v,
	}
}

func (s *NotesService) ListNotes(_ *notesv1.ListNotesRequest, srv notesv1.NotesService_ListNotesServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, note := range s.notes {
		err := srv.Send(&notesv1.ListNotesResponse{Note: note})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *NotesService) AddNote(_ context.Context, req *notesv1.AddNoteRequest) (*notesv1.AddNoteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	note := &notesv1.Note{
		Id:       uuid.Must(uuid.NewV4()).String(),
		Title:    req.Title,
		Archived: false,
	}

	if err := s.validator.Validate(note); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	s.notes = append(s.notes, note)

	return &notesv1.AddNoteResponse{
		Note: note,
	}, nil
}
