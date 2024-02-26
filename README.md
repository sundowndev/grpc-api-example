# grpc-api-example

This repository demonstrates how to use gRPC with Go and various plugins. It's a simple API to manage notes. It uses grpc-gateway to handle HTTP requests and openapi to define API specs.

## Installation

This project uses go1.20.

```
$ make install-tools
$ make generate
```

### Run the server

```
$ make build
$ ./bin/server
```