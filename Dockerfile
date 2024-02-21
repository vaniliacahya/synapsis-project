# Use a minimal base image for Go applications
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application binary from the builder stage into the final container
COPY --from=builder /app/main .

# Expose the port that the Go application listens on
EXPOSE 8080

# Command to run the Go application binary
CMD ["./main"]

