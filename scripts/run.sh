#!/bin/bash

export MONGO_URI=mongodb://config-server:config-server-password@localhost:27017 
export PORT=8081

go run cmd/grpc/grpc.go
