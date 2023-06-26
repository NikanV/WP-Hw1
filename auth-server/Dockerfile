FROM golang:latest

WORKDIR /WP-Hw1

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./auth-server/server.go ./
COPY ./proto/ ./proto/
COPY ./tools/ ./tools/

RUN go build -o ./auth-server

EXPOSE 5052

CMD ["./auth-server"]