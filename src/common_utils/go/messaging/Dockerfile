# Use the official Go image as the base image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies - this happens during the build
RUN go mod download

# Copy the rest of your application code
COPY . .

# Expose any ports your application might need (optional for this test)
# EXPOSE 8080

# Command to run when the container starts
# This executes your Go tests with verbose output
# The TEST_KAFKA_BROKER environment variable will be provided by docker-compose
CMD ["go", "test", "-v", "./..."]
