#!/bin/bash

# this scripts builds and tags with docker

# it takes one argument, the tag version

TAG=latest

if [ $# -eq 1 ]; then
    TAG=$1
fi

docker build -t dekuyo/fingo-auth:$TAG -f _deployments/Dockerfile.auth --target production .
docker build -t dekuyo/fingo-token:$TAG -f _deployments/Dockerfile.token --target production .
docker build -t dekuyo/fingo-contact:$TAG -f _deployments/Dockerfile.contact --target production .
docker build -t dekuyo/fingo-wallet:$TAG -f _deployments/Dockerfile.wallet --target production .

docker push dekuyo/fingo-auth:$TAG
docker push dekuyo/fingo-token:$TAG
docker push dekuyo/fingo-contact:$TAG
docker push dekuyo/fingo-wallet:$TAG
