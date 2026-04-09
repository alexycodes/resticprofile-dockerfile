FROM debian:13.4-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    vim

COPY --from=golang:1.26-alpine /usr/local/go/ /usr/local/go/
COPY --from=golangci/golangci-lint:v2.11-alpine /usr/bin/golangci-lint /usr/local/go/bin/golangci-lint

ENV PATH="$PATH:/usr/local/go/bin:/root/go/bin"

RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest
