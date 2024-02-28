package server

import (
	"context"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"github.com/gofrs/uuid"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	s.notes = append(s.notes, note)

	return &notesv1.AddNoteResponse{Note: note}, nil
}

func (s *NotesService) EditNote(_ context.Context, req *notesv1.EditNoteRequest) (*notesv1.EditNoteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	editedNote := &notesv1.Note{
		Id:       req.Note.Id,
		Title:    req.Note.Title,
		Archived: req.Note.Archived,
	}

	if err := s.validator.Validate(editedNote); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	var found bool
	for i := range s.notes {
		if s.notes[i].Id == req.Note.Id {
			s.notes[i] = editedNote
			found = true
			break
		}
	}

	if !found {
		return nil, status.New(codes.NotFound, fmt.Sprintf("couldn't find note with id: %s", req.Note.Id)).Err()
	}

	return &notesv1.EditNoteResponse{Note: editedNote}, nil
}
