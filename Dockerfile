FROM golang:1.11.1-alpine3.8

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN apk update
RUN apk add --no-cache git

RUN go mod download \
    go build -o main /

CMD ["/app/main"]