# Use the official Go image for building the application
FROM golang:1.23.4 AS builder

# Set the working directory inside the container
WORKDIR /

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

# Use a minimal base image for the final container
FROM gcr.io/distroless/base-debian10

# Set the working directory in the container
WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder main .

# Set the entrypoint command
CMD ["./main"]