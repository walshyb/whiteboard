FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache make protobuf-dev

# Install Go Protobuf Plugins
# These are the executables that protoc calls (protoc-gen-go)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go.mod and go.sum first to allow caching dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go mod tidy

COPY . .

RUN make proto_go

RUN go build -mod=mod -ldflags='-s -w' -o /bin/server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/server /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]
