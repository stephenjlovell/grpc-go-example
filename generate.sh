#!/bin/bash

protoc --proto_path=api/ --go_out=plugins=grpc:api/go/pkg --ruby_out=api/ruby api/*.proto
