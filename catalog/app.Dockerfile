# Build stage
FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

COPY vendor vendor
COPY catalog catalog

# Build the application
RUN GO111MODULE=on go build -mod vendor -o main ./catalog/cmd/catalog

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=build /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]