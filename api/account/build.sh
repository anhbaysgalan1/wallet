#!/bin/bash

# 编译google.api
#protoc -I=. \
#  --go_opt=paths=source_relative \
  #  -I=$GOPATH/pkg/mod/github.com/protocolbuffers/protobuf@v3.19.1+incompatible/src \
#  --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. \
#  ../google/api/*.proto


## 编译hello_http.proto
protoc -I=. \
  -I=$GOPATH/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20211108224529-691a18b1c8af \
  -I=$GOPATH/pkg/mod/protobuf \
  --go_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/leaf-rain/wallet/api/google/api:. \
  ./*.proto
#
## 编译hello_http.proto gateway
#protoc -I=. \
#  -I=$GOPATH/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20211108224529-691a18b1c8af \
#  -go_opt=paths=source_relative
#  --grpc-gateway_out=logtostderr=true:. \
#  ./accout.proto
