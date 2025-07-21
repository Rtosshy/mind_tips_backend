FROM golang:latest AS builder

WORKDIR /app

SHELL ["/bin/bash", "-c"]

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    bash \
    git && \
    rm -rf /var/lib/apt/lists/*

ARG GOLANGCILINT_VERSION=v2.1.6
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCILINT_VERSION}

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

ENV PATH="/usr/local/go/bin:$(go env GOPATH)/bin:/root/go/bin:${PATH}"

CMD ["air", "-c", ".air.toml"]

