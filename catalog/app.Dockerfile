# Build stage
FROM golang:1.23.4-alpine3.21 AS build
RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

COPY catalog catalog

# Build the application
RUN GO111MODULE=on go build -o main ./catalog/cmd/catalog

# Final stage
FROM alpine:3.21

WORKDIR /app

# Copy binary from builder
COPY --from=build /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]