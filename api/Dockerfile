FROM golang:1.25-alpine AS build

RUN adduser -D -u 1000 nyaa
USER nyaa

WORKDIR /app

# Restart on *.go change
RUN go install github.com/mitranim/gow@latest

# Debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install ./cmd/goose

ARG COMMIT_HASH
ENV COMMIT_HASH=$COMMIT_HASH

CMD goose up && gow run -ldflags="-X codeberg.org/jessienyan/booruview.AppVersion=$COMMIT_HASH" cmd/server/main.go
