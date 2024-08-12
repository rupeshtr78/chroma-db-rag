#!/bin/bash

mkdir -p ./internal/grpc/generated
protoc --go_out=./internal/grpc/generated --go-grpc_out=./internal/grpc/generated ./internal/grpc/tei.proto

# split the above command into two commands
# option go_package = "internal/grpc/generated;generated";
protoc --go_out=. -I ./internal/grpc ./internal/grpc/tei.proto
protoc --go-grpc_out=. -I ./internal/grpc ./internal/grpc/tei.proto

# Embedding with curl
curl 127.0.0.1:8080/embed \
    -X POST \
    -d '{"inputs":"What is Deep Learning?"}' \
    -H 'Content-Type: application/json'

# Embedding with grpcurl
grpcurl -d '{"inputs": "What is Deep Learning"}' -plaintext 10.0.0.213:50082 tei.v1.Embed/Embed

# Reranking with curl
curl 127.0.0.1:8080/rerank \
    -X POST \
    -d '{"query":"What is Deep Learning?", "texts": ["Deep Learning is not...", "Deep learning is..."]}' \
    -H 'Content-Type: application/json'

# Reranking with grpcurl
grpcurl -d '{"query":"What is Deep Learning?", "texts": ["Deep Learning is not...", "Deep learning is..."], "return_text": true}' \
    -plaintext 10.0.0.213:50083 \
    tei.v1.Rerank/Rerank