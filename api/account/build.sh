#!/bin/bash

protoc -I=. --go_out=plugins=grpc:. *.srv.proto

#protoc -I=. --go_out=../../internal/tp_wallet/application/dto --go_opt=paths=source_relative \
#  *.msg.proto
#
#protoc -I=. --go-grpc_out=../../internal/tp_wallet/application/api --go-grpc_opt=paths=source_relative \
#  *.srv.proto