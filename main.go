package main

import (
	"errors"
	"github.com/sundowndev/grpc-api-example/server"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
	grpclog.SetLoggerV2(log)

	addr := "0.0.0.0:10000"                            // TODO: handle this with CLI flags
	srv := server.NewServer(insecure.NewCredentials()) // TODO: Replace with your own certificate!

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		if err := srv.Listen(addr); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Fatalf("listen: %s\n", err)
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

	if err := srv.Close(); err != nil {
		log.Fatal(err)
	}

	log.Info("graceful shut down succeeded")
}
