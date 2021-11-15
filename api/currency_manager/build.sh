#!/bin/bash

protoc -I=. --go_out=. --go_opt=paths=import \
    --go-grpc_out=. --go-grpc_opt=paths=import \
    ./*.srv.proto

