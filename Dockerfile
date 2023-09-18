FROM golang:1.21-alpine AS builder

RUN apk mkdir /cmd
ADD . /cmd/
WORKDIR /cmd
RUN go build -o main
CMD ["/cmd/main"]
