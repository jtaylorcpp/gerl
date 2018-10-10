#!/bin/bash
GRPC=/core/
GRPCDIR=$(pwd)$GRPC
protoc -I=$GRPCDIR --go_out=plugins=grpc:$GRPCDIR $GRPCDIR/grpc.proto
