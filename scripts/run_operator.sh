#!/bin/bash

WORKSPACE=${WORKSPACE:-$HOME/tf-operator}
kubectl apply -k ${WORKSPACE}/deploy/kustomize/operator/latest/
kubectl apply -k ${WORKSPACE}/deploy/kustomize/contrail/1node/latest/
