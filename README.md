# grpc-api-example

This repository demonstrates how to use gRPC with Go and various plugins. It's a simple API to manage notes. It uses grpc-gateway to handle HTTP requests and openapi to define API specs.

## Demonstrated features

- buf.build framework
- grpc gateway + openapi specs and swagger UI
- graceful server shut down (designed for orchestration systems such as k8s)
- gRPC [health check protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md#grpc-health-checking-protocol).
- rpc message validation

⚠️ For demo purposes, encryption wasn't properly handled in this project.

## Installation

This project uses go1.20.

```
$ make install-tools
$ make build
```

### Run the server

```
$ ./bin/server -h
Usage of ./bin/server:
  -grpc-server-endpoint string
    	gRPC server endpoint (default "localhost:9090")
  -http-server-endpoint string
    	HTTP server endpoint (default "localhost:8000")
```

#### Health checking

```
$ ./bin/
```
