# --- Stage 1: The Builder ---
FROM golang:tip-alpine3.21 AS builder

# Ensure the ca-certificates package is installed in the builder.
# This gives us the root certificates we need to copy later. This command
# works on any Alpine-based image.
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container.
WORKDIR /app

# Copy the dependency files first for caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of our application source code into the container.
COPY . .

# Build the Go application.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/server ./cmd/api/main.go


# --- Stage 2: The Final Image ---
# We use 'alpine' as our base. It's a tiny, secure Linux
# distribution that provides a minimal environment, including networking tools.
FROM alpine:latest

# Set the working directory.
WORKDIR /app

# THE FIX: Copy the trusted CA certificates from our 'builder' stage.
# This gives our application the "phonebook" it needs to verify HTTPS connections.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the compiled binary from the 'builder' stage.
COPY --from=builder /bin/server .

# Copy the prompts directory from the 'builder' stage.
COPY --from=builder /app/prompts ./prompts

# Expose port 8080.
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["./server"]
