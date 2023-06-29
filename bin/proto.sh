#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 <relative path> <filename>"
  exit 1
fi

type=$1
filename2=$2

protoc --go_out=. --go_opt=paths=$type --go-grpc_out=. --go-grpc_opt=paths=$type $filename2