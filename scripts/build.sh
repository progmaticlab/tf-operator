#!/bin/bash -e

type go >/dev/null 2>&1 || {
  export PATH=$PATH:/usr/local/go/bin
}

export WORKSPACE=${WORKSPACE:-$HOME}
export CONTRAIL_REPOSITORY=${CONTRAIL_REPOSITORY:-"localhost:5000"}
export CONTRAIL_CONTAINER_TAG=${CONTRAIL_CONTAINER_TAG:-"latest"}
export CGO_ENABLED=1

target=${CONTRAIL_REPOSITORY}/${CONTRAIL_CONTAINER_TAG}
pushd ${WORKSPACE}/tf-operator

operator-sdk build $target
docker push $target

popd
