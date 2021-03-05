#!/bin/bash

protoc --proto_path=api/ --go-grpc_out=api/go/pkg --go_out=api/go/pkg --ruby_out=api/ruby api/*.proto
