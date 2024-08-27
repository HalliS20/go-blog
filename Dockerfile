FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY app/ .

# Install sqlite3 and its development libraries
RUN apt-get update && apt-get install -y gcc

# for .env file
# RUN export $(cat .env | xargs)

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o main

EXPOSE 8080

# Set the entry point for the container
CMD ["/app/main"]


