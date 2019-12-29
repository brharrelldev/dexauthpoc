proto-demo:
	protoc -I. -I$(GOPATH)/src  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapi --go_out=plugins=grpc:. api/*.proto

.PHONY: build-client build-server build-all

build-auth:
	go build -o auth-server authserver/auth.go

build-server:
	go build -o grpc-server server/server.go

build-client:
	go build -o grpc-client client/client.go


build-all: build-auth build-client build-server

launch-demo:
	itermocil oauth-demo



