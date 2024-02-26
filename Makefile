install-tools:
	go install gotest.tools/gotestsum@v1.11.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/bufbuild/buf/cmd/buf@v1.29.0

build: generate
	buf build
	go build -o bin/server .

generate:
	buf generate
	go generate ./...

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