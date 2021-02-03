#!/bin/bash -e

WORKSPACE=${WORKSPACE:-$HOME/tf-operator}

OPERATOR_VERSION=${OPERATOR_VERSION:-'latest'}
CONTRAIL_VERSION=${CONTRAIL_VERSION:-'tf-ci-nightly'}
CONTRAIL_HA=${CONTRAIL_HA:-'1node'}

kubectl apply -k ${WORKSPACE}/deploy/kustomize/operator/$OPERATOR_VERSION/
while [[ ! $(kubectl wait crds --for=condition=Established --timeout=2m managers.contrail.juniper.net) ]]
do
  sleep 2s
done

kubectl apply -k ${WORKSPACE}/deploy/kustomize/contrail/$CONTRAIL_HA/$CONTRAIL_VERSION/
