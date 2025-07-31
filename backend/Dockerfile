# --- Stage 1: The Builder ---
# We start with our Alpine-based Go image.
FROM golang:tip-alpine3.21 AS builder

# NEW: Install the C compiler and build tools (gcc, make, etc.).
# The 'build-base' package on Alpine contains all the necessary tools for CGO.
# We also still need the ca-certificates for HTTPS calls.
RUN apk add --no-cache build-base ca-certificates

# Set the working directory inside the container.
WORKDIR /app

# Copy the dependency files first for caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of our application source code into the container.
COPY . .

# Build the Go application.
# REMOVED: We have removed 'CGO_ENABLED=0' to allow our program to use
# the go-sqlite3 C bindings.
RUN go build -ldflags="-w -s" -o /bin/server ./cmd/api/main.go


# --- Stage 2: The Final Image ---
# We continue to use 'alpine' as our base. This is now CRITICAL because
# our compiled program is no longer fully static. It needs the standard
# C libraries (like musl libc) that Alpine provides to run correctly.
FROM alpine:latest

# Set the working directory.
WORKDIR /app

# Copy the trusted CA certificates from our 'builder' stage.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the compiled binary from the 'builder' stage.
COPY --from=builder /bin/server .

# Copy the prompts directory from the 'builder' stage.
COPY --from=builder /app/prompts ./prompts

# Expose port 8080.
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["./server"]
