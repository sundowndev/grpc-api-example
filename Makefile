BUF_VERSION:=v1.29.0
SWAGGER_UI_VERSION:=v5.11.8

install-tools:
	go install gotest.tools/gotestsum@v1.11.0
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1

build: generate
	buf build
	go build -o bin/server .

generate:
	find ./proto -type f -name "*.go" -delete # Delete existing Go generated files
	rm -rf ./gen/openapiv2/* # Delete existing openapi specs
	buf generate
	go generate ./...
	SWAGGER_UI_VERSION=$(SWAGGER_UI_VERSION) ./scripts/generate-swagger-ui.sh

lint:
	buf lint

fmt:
	buf format
	go fmt ./...

go.mod:
	go mod tidy
	go mod verify
go.sum: go.mod

test:
	gotestsum --format testname --junitfile unit-tests.xml -- -mod=readonly -race -coverprofile=./c.out -covermode=atomic -coverpkg=.,./... ./...

coverage: test
	go tool cover -func=cover.out

mocks:
	rm -rf mocks/*
	mockery --all

clean:
	go clean
	rm -f bin/*