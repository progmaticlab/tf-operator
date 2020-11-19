#!/bin/bash

WORKSPACE=${WORKSPACE:-tf-operator}
kubectl apply -k ${WORKSPACE}/deploy/kustomize/operator/latest/
kubectl apply -k deploy/kustomize/contrail/1node/latest/
