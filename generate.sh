#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet/greetpb/greet.proto
#protoc greet/greetpb/greet.proto --go-grpc_out=grpc:.