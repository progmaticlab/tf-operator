#!/bin/bash

operator-sdk build localhost:5000/contrail-operator:latest
docker push localhost:5000/contrail-operator:latest

# build CRDS container

docker build --tag localhost:5000/crdsloader:latest deploy/crds
localhost:5000/crdsloader:latest
