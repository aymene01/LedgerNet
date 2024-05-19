.PHONY: proto

build:
	@go build ./cmd/main.go -o bin/blocker

run: build
	@./bin/blocker

clean:
	@rm -rf bin

test:
	@go test -v ./...

proto:
	protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.
