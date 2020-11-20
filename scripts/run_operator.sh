#!/bin/bash

WORKSPACE=${WORKSPACE:-$HOME/tf-operator}
kubectl apply -k ${WORKSPACE}/deploy/kustomize/operator/latest/
while [[ ! $(kubectl wait crds --for=condition=Established --timeout=2m managers.contrail.juniper.net) ]]
do
  sleep 2s
done

kubectl apply -k ${WORKSPACE}/deploy/kustomize/contrail/1node/latest/
