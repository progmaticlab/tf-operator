#!/bin/bash

operator-sdk build localhost:5000/tf-operator:latest
docker push localhost:5000/tf-operator:latest

# build CRDS container

docker build --tag localhost:5000/crdsloader:latest deploy/crds
docker push localhost:5000/crdsloader:latest
