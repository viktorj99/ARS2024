# syntax=docker/dockerfile:1

# Use an official Golang runtime as a parent image
FROM golang:1.22.1 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./
# Download any necessary dependencies
RUN go mod download

# Copy the source code into the container's working directory
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /projekat

# Use a small Alpine Linux image for the final stage
FROM alpine:latest
WORKDIR /root/

# Install ca-certificates, required for many Go applications
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /projekat .

# Make port 8080 available to the outside world
EXPOSE 8080

# Run the binary
CMD ["./projekat"]