#!/bin/bash

PROTOS="${GOPATH}/src/github.com/bbengfort/otterdb/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

MODULE="github.com/bbengfort/otterdb/pkg/replica/raft/v1"
APIMOD="github.com/bbengfort/otterdb/pkg/replica/raft/v1;raft"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 --go-grpc_out=./v1 \
    --go_opt=module=${MODULE} \
    --go-grpc_opt=module=${MODULE} \
    --go_opt=Mraft/v1/raft.proto="${APIMOD}" \
    --go-grpc_opt=Mraft/v1/raft.proto="${APIMOD}" \
    raft/v1/raft.proto