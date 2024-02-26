BUF_VERSION:=v1.29.0

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

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

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod

test:
	$(GOTEST) --format testname --junitfile unit-tests.xml -- -mod=readonly -race -coverprofile=./c.out -covermode=atomic -coverpkg=.,./... ./...

coverage: test
	$(GOTOOL) cover -func=cover.out

mocks:
	rm -rf mocks/*
	mockery --all
