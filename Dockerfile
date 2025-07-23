# Use the official Golang image as a builder
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# Use a minimal base image for the final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server .

# Expose the port the application listens on
EXPOSE 8080

# Run the application
CMD ["/app/server"]
