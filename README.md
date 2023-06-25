# WP-Hw1
In order to use HTTPS with nginx you can use the [MKCERT](https://github.com/FiloSottile/mkcert) repository

Use the commands below from the root directory to setup your environment for the project

Compiling the `.proto` files :
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

Compiling the `swagger` files :
```bash
swag init --parseDependency --parseInternal --parseDepth 1 -g /gateway-server/main.go
```

Building images from servers `Dockerfiles` :
```bash
docker build --tag gateway-server -f ./gateway-server/Dockerfile .
docker build --tag biz-server -f ./biz-server/Dockerfile .
docker build --tag auth-server -f ./auth-server/Dockerfile .
```

Creating a network in docker :
```bash
docker network create serversdb
```

Use this commands in `docker-serversdb` and `docker-nginx` directories in order to run the containers :

```bash
docker compose up -d
```