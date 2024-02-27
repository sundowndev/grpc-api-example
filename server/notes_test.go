package server

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"net"
	"testing"
)

func newTestServer() (*Server, *grpc.ClientConn, error) {
	srv, err := NewServer(insecure.NewCredentials())
	if err != nil {
		return nil, nil, err
	}
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)
	srv.listener = lis
	go srv.grpcSrv.Serve(lis)

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		_ = srv.Close()
		return nil, nil, err
	}

	return srv, conn, nil
}

func TestNotesService_ListNotes(t *testing.T) {
	testcases := []struct {
		name    string
		notes   []*notesv1.Note
		wantErr string
	}{
		{
			name:  "test with no notes",
			notes: []*notesv1.Note{},
		},
		{
			name: "test with few notes",
			notes: []*notesv1.Note{
				{Title: "test note 1"},
				{Title: "note 2"},
				{Title: "note 3"},
			},
		},
		{
			name: "test with min_len validation error",
			notes: []*notesv1.Note{
				{Title: ""},
			},
			wantErr: "validation failed: validation error:\n - title: value length must be at least 1 characters [string.min_len]",
		},
		{
			name: "test with max_len validation error",
			notes: []*notesv1.Note{
				{Title: "this is a super long note title that can trigger a validation error"},
			},
			wantErr: "validation failed: validation error:\n - title: value length must be at most 50 characters [string.max_len]",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			srv, conn, err := newTestServer()
			if err != nil {
				t.Fatal(err)
			}
			client := notesv1.NewNotesServiceClient(conn)
			defer srv.Close()
			defer conn.Close()

			for _, note := range tt.notes {
				_, err := client.AddNote(context.Background(), &notesv1.AddNoteRequest{Title: note.Title})
				if err != nil && tt.wantErr == "" {
					t.Fatal(err)
				} else {
					// Retrieve the status from the rpc error
					s, ok := status.FromError(err)
					if ok {
						assert.Equal(t, tt.wantErr, s.Message())
					} else {
						log.Fatal(err)
					}
					return
				}
			}

			res, err := client.ListNotes(context.Background(), &notesv1.ListNotesRequest{})
			if err != nil {
				t.Fatal(err)
			}

			notes := make([]*notesv1.Note, 0)

			for {
				recv, err := res.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				notes = append(notes, recv.Note)
			}

			assert.Len(t, notes, len(tt.notes))
		})
	}
}

func TestNotesService_AddNote(t *testing.T) {
	srv, conn, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	client := notesv1.NewNotesServiceClient(conn)
	defer srv.Close()
	defer conn.Close()

	res, err := client.AddNote(context.Background(), &notesv1.AddNoteRequest{Title: "test note"})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, res)
	assert.Equal(t, "test note", res.Note.Title)
	assert.Equal(t, false, res.Note.Archived)
}

func TestNotesService_EditNote(t *testing.T) {
	srv, conn, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	client := notesv1.NewNotesServiceClient(conn)
	defer srv.Close()
	defer conn.Close()

	res, err := client.AddNote(context.Background(), &notesv1.AddNoteRequest{Title: "test note"})
	if err != nil {
		t.Fatal(err)
	}

	res.Note.Title = "january groceries"
	res.Note.Archived = true

	res2, err := client.EditNote(context.Background(), &notesv1.EditNoteRequest{Note: res.Note})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "january groceries", res2.Note.Title)
	assert.True(t, res2.Note.Archived)
}
