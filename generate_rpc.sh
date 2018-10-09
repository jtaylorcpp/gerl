#!/bin/bash
GRPC=/core/includes/grpc
GRPCDIR=$(pwd)$GRPC
protoc -I=$GRPCDIR --go_out=plugins=grpc:$GRPCDIR $GRPCDIR/grpc.proto
