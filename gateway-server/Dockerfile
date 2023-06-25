FROM golang:latest

WORKDIR /WP-Hw1

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./gateway-server/main.go ./
COPY ./proto/ ./proto/
COPY ./tools/ ./tools/
COPY ./docs/ ./docs/

RUN go build -o ./gateway-server

EXPOSE 6443

CMD ["./gateway-server"]