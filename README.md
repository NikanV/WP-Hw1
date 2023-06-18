# WP-Hw1
In order to use HTTPS with nginx you can use the [MKCERT](https://github.com/FiloSottile/mkcert) repository

Use the command below in the wanted server directory in order to compile `.proto` files :
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    "folder name"/"file name".proto
```

use this to save swagger changes
```bash
swag init --parseDependency --parseInternal --parseDepth 1 -g /gateway-server/main.go
```