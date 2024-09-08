# Use a base image that includes Caddy
FROM caddy:2 AS caddy

# Install necessary packages for building Go app
RUN apk add --no-cache git build-base

# Set up Caddy configuration
COPY Caddyfile /etc/caddy/Caddyfile


FROM golang:1.22-alpine AS builder

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
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk add --no-cache caddy

WORKDIR /root/

COPY --from=caddy /etc/caddy/Caddyfile /etc/caddy/Caddyfile

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY app/templates /root/templates
COPY app/public /root/public

# Expose port 8080
EXPOSE 80 443

# Command to run both the Go application and Caddy
CMD ["sh", "-c", "caddy run --config /etc/caddy/Caddyfile & ./main"]