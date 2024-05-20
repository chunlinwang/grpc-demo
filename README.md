# This is a grpc demo project

### How to install protoc?
[https://google.github.io/proto-lens/installing-protoc.html](https://google.github.io/proto-lens/installing-protoc.html)

## proto buffers file in the notification dir.

### How to compile proto buffers

```shell
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  notification/notification.proto
```

### start the server
```shell
go run notification_server/main.go
```

### start the client
```shell
go run notification_client/main.go
```

## Blog
https://medium.com/@kazami0083/grpc-vs-restful-api-49753545e3cf

## requirement
- Go > 1.20
- libprotoc 26.1

## Author
* [@Chunlin Wang](https://www.linkedin.com/in/chunlin-wang-b606b159/)