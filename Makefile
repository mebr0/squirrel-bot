.PHONY:
.SILENT:
.DEFAULT_GOAL := build

lint:
	golangci-lint run

test:
	go test --short -coverprofile=cover.out -v ./...

cover: test
	go tool cover -func=cover.out

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

run: build
	./app

fmt:
	gofmt -s -w .