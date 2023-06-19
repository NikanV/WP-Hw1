FROM golang:latest

WORKDIR /WP-Hw1

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./biz-server/server.go ./
COPY ./proto/ ./proto/

RUN go build -o ./biz-server

EXPOSE 5062

CMD ["./biz-server"]