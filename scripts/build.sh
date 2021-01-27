#!/bin/bash -e

type go >/dev/null 2>&1 || {
  export PATH=$PATH:/usr/local/go/bin
}
CONTAINER_REGISTRY=${CONTAINER_REGISTRY:-"localhost:5000"}
CONTRAIL_CONTAINER_TAG=${CONTRAIL_CONTAINER_TAG:-"latest"}

export  CGO_ENABLED=1
operator-sdk build  ${CONTAINER_REGISTRY}/tf-operator:${CONTRAIL_CONTAINER_TAG}
sudo docker push ${CONTAINER_REGISTRY}/tf-operator:${CONTRAIL_CONTAINER_TAG}

# build CRDS container
sudo docker build --tag ${CONTAINER_REGISTRY}/tf-crdsloader:${CONTRAIL_CONTAINER_TAG} deploy/crds
sudo docker push ${CONTAINER_REGISTRY}/tf-crdsloader:${CONTRAIL_CONTAINER_TAG}
