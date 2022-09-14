#!/usr/bin/env bash

# Uses https://github.com/grpc/grpc-go and https://pkg.go.dev/github.com/golang/protobuf/protoc-gen-go as well as https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc for Go
# Uses https://www.npmjs.com/package/@grpc/grpc-js and https://github.com/thesayyn/protoc-gen-ts for Node.js + TypeScript

# Installation prerequisites:
# brew install protobuf
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
# export PATH="$PATH:$(go env GOPATH)/bin"
# pnpm install

# Generate protobuf + gRPC for Go
protoc \
  --go_out=./protobuf --go_opt=paths=source_relative \
  --go-grpc_out=./protobuf --go-grpc_opt=paths=source_relative \
  sdk.proto

# store TypeScript code in sdk/atlas-sdk-ts/src/protobuf
protoc --ts_out=./sdk/atlas-sdk-ts/src --ts_opt=unary_rpc_promise=true --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts sdk.proto
