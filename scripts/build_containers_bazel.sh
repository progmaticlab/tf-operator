#!/bin/bash

WORKSPACE=${WORKSPACE:-$HOME/tf-operator}
if [[ ! -d ${HOME}/env ]]; then
  python3 -m venv ${HOME}/env
fi

source ${HOME}/env/bin/activate
pushd $WORKSPACE
bazel run //cmd/crdsloader:contrail-operator-crdsloader-push-local
bazel run //cmd/manager:contrail-operator-push-local
bazel run //contrail-provisioner:contrail-provisioner-push-local
bazel run //statusmonitor:contrail-statusmonitor-push-local
popd
