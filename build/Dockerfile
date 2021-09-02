# Main dockerfile for building images on unix systems (build can be done manually)
#
# Build binary:
# go build -ldflags="-w -s" -v -o app ./cmd/main.go
#
# Build image:
# docker build -f build/Dockerfile -t ghcr.io/mebr0/squirrel-bot .
#
# Run container by passing env variables by flag -e (list of env variables in README.md):
# docker run ghcr.io/mebr0/squirrel-bot:latest
#
# Transfer artifacts to clean alpine container
FROM alpine:latest

WORKDIR /root/

COPY . .

CMD ["./app"]
