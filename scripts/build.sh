#!/bin/bash -e
CONTAINER_REGISTRY=${CONTAINER_REGISTRY:-"localhost:5000"}

export  CGO_ENABLED=1
operator-sdk build  ${CONTAINER_REGISTRY}/tf-operator:latest
sudo docker push ${CONTAINER_REGISTRY}/tf-operator:latest

# build CRDS container
sudo docker build --tag ${CONTAINER_REGISTRY}/tf-crdsloader:latest deploy/crds
sudo docker push ${CONTAINER_REGISTRY}/tf-crdsloader:latest
