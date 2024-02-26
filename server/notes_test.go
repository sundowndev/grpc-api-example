package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestNotesService_ListNotes(t *testing.T) {
	addr := "0.0.0.0:10000"
	srv := NewServer()
	go srv.Listen(addr)
	defer srv.Close()

	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := notesv1.NewNotesServiceClient(conn)
	res, err := client.ListNotes(context.TODO(), &notesv1.ListNotesRequest{})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, res)

	//err = res.RecvMsg()
	//if err != nil && !errors.Is(err, io.EOF) {
	//	t.Fatal(err)
	//}
}

func TestNotesService_AddNote(t *testing.T) {
	addr := "0.0.0.0:10000"
	srv := NewServer()
	go srv.Listen(addr)
	defer srv.Close()

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := notesv1.NewNotesServiceClient(conn)
	res, err := client.AddNote(context.TODO(), &notesv1.AddNoteRequest{Title: "test note"})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, res)
}
