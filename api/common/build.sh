#!/bin/bash

protoc -I=. --go_out=../../internal/common/ *.common.proto
