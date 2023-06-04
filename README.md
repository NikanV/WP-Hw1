# WP-Hw1
In order to use HTTPS with nginx you can use the [MKCERT](https://github.com/FiloSottile/mkcert) repository

Use the command below in the root folder in order to compile `.proto` files :
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    "folder name"/"file name".proto
```
