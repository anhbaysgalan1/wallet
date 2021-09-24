#!/bin/bash

protoc -I=. --go_out=../common/ ./currency.proto ./net.proto