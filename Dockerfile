# Build stage for Go application and Protobuf tools
FROM golang:1.25-bookworm AS builder

# Install protobuf compiler and other necessary tools
RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    make \
    dnsutils \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

# Install Go plugins for protoc and other tools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

ENV PATH="/go/bin:${PATH}"
RUN ln -s /go/bin/goose /usr/local/bin/goose

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Generate protobuf code (optional but good for consistency)
RUN make gen-go

# Make entrypoint script executable
RUN chmod +x entrypoint.sh

# The application runs using go run, so we keep the Go environment
ENTRYPOINT ["./entrypoint.sh"]
