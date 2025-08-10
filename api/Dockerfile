FROM golang:1.24-alpine AS build

WORKDIR /app

# Restart on *.go change
RUN go install github.com/mitranim/gow@latest

# Debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG COMMIT_HASH
ENV COMMIT_HASH=$COMMIT_HASH

CMD gow run -ldflags="-X github.com/jessienyan/booruview.AppVersion=$COMMIT_HASH" cmd/main.go
