FROM golang:alpine

WORKDIR /app
COPY go.mod ./
COPY . .

RUN go build -o main ./cmd/main.go

CMD ["./main"]