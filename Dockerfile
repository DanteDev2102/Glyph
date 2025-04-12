FROM golang:1.24.1-bookworm

WORKDIR /app

COPY . .

RUN go install golang.org/x/lint/golint@latest

RUN go install golang.org/x/tools/cmd/goimports@latest

RUN go install github.com/air-verse/air@latest