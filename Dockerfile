FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY app/ .

# Install gcc and its dependencies
RUN apk add --no-cache gcc musl-dev

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Expose both TCP and UDP on port 443
EXPOSE 443/tcp
EXPOSE 443/udp

# Expose port 8080 for HTTP (if needed)
EXPOSE 8080/tcp


# Set the entry point for the container
CMD ["/app/main"]
