.PHONY:
.SILENT:
.DEFAULT_GOAL := compile

lint:
	golangci-lint run

test:
	go test --short -coverprofile=cover.out -v ./...

cover: test
	go tool cover -func=cover.out

compile:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

run: compile
	./app

fmt:
	gofmt -s -w .