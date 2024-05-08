# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .


RUN go mod download

COPY . .

RUN go build -o server .

# Final stage
FROM redis:latest

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]