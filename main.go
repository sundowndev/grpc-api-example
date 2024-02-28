package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sundowndev/grpc-api-example/gen"
	notesv1 "github.com/sundowndev/grpc-api-example/proto/notes/v1"
	"github.com/sundowndev/grpc-api-example/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"io"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	httpServerEndpoint = flag.String("http-server-endpoint", "localhost:8000", "HTTP server endpoint")
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func getOpenAPIHandler() (http.Handler, error) {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return nil, err
	}

	// Use subdirectory in embedded files
	subFS, err := fs.Sub(gen.OpenAPI, "openapiv2")
	if err != nil {
		return nil, fmt.Errorf("couldn't create sub filesystem: %v", err)
	}
	return http.FileServer(http.FS(subFS)), nil
}

func httpServer(ctx context.Context, addr string) error {
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: Replace with your own certificate!
	}

	// Register services
	err := notesv1.RegisterNotesServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	oa, err := getOpenAPIHandler()
	if err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr: *httpServerEndpoint,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				mux.ServeHTTP(w, r)
				return
			}
			oa.ServeHTTP(w, r)
		}),
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return gwServer.ListenAndServe()
}

func main() {
	flag.Parse()

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
	grpclog.SetLoggerV2(log)

	srv, err := server.NewServer(insecure.NewCredentials()) // TODO: Replace with your own certificate!
	if err != nil {
		log.Fatal("failed to initialize server: %v", err)
	}

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", *grpcServerEndpoint)
	go func() {
		if err := srv.Listen(*grpcServerEndpoint); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Serve HTTP gateway
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := httpServer(ctx, *httpServerEndpoint); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	if err := srv.Close(); err != nil {
		log.Fatal(err)
	}

	log.Info("graceful shut down succeeded")
}
