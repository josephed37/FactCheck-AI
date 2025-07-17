# We start with a full Go environment to build our application.
# We use a specific version to ensure consistent builds.
FROM golang:1.22-focal AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the dependency files first.
# This is a Docker caching optimization. As long as these files don't change,
# Docker won't re-download the dependencies on subsequent builds, making them much faster.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of our application source code into the container.
COPY . .

# Build the Go application.
# - CGO_ENABLED=0 disables C bindings, creating a purely static binary.
# - -o /bin/server builds the output file and names it 'server'.
# - -ldflags="-w -s" strips debugging information, making the binary smaller.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/server ./cmd/api/main.go


# We start from a minimal, empty base image called "scratch".
# This is the most secure base image because it contains nothing but our application.
FROM scratch

# Set the working directory.
WORKDIR /app

# Copy only the compiled binary from the 'builder' stage.
# We are not copying any of the source code or build tools, which keeps our
# final image incredibly small and secure.
COPY --from=builder /bin/server .

# Copy the prompts directory from the 'builder' stage.
# This is necessary because our app reads the prompt file at runtime.
COPY --from=builder /app/prompts ./prompts

# Expose port 8080. This tells Docker that our application inside the
# container will be listening on this port.
EXPOSE 8080

# The command to run when the container starts.
# This executes our compiled application.
ENTRYPOINT ["./server"]
