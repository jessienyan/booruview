FROM golang:1.24-alpine AS build

WORKDIR /app

# Restart on *.go change
RUN go install github.com/mitranim/gow@latest

# Debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY get_goose.sh .
RUN ./get_goose.sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG COMMIT_HASH
ENV COMMIT_HASH=$COMMIT_HASH

CMD go run cmd/createdb/main.go && goose up && gow run -ldflags="-X codeberg.org/jessienyan/booruview.AppVersion=$COMMIT_HASH" cmd/server/main.go
