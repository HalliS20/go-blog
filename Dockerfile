FROM golang:1.20

WORKDIR ./

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev
RUN CGO_ENABLED=1 GOOS=linux go build -v -o main .

CMD ["./main"]
