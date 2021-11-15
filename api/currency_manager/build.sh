#!/bin/bash

protoc -I=. --go_out=plugins=grpc:. *.srv.proto

