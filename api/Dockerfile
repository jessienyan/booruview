FROM golang:1.24-alpine AS build

ENV VITE_COMMIT_SHA=$VITE_COMMIT_SHA

WORKDIR /app

# Restart on *.go change
RUN go install github.com/mitranim/gow@latest

# Debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD gow run -ldflags="-X github.com/jessienyan/booruview.AppVersion=$VITE_COMMIT_SHA" cmd/main.go
