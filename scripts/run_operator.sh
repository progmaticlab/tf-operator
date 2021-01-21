#!/bin/bash -e

WORKSPACE=${WORKSPACE:-$HOME/tf-operator}

# multinode setup with registry on one node
count=$(echo "$CONTROLLER_NODES $AGENT_NODES" | tr " " "\n" | sort -u |  tr "\n" " " | awk -F ' ' '{print NF}' )
if [[ $count > 1 ]] ; then
  registry_ip=$(hostname -i)
  sed -i "s/localhost/$registry_ip/g" ${WORKSPACE}/deploy/kustomize/operator/latest/*
  sed -i "s/localhost/$registry_ip/g" ${WORKSPACE}/deploy/kustomize/base/operator/*
  sed -i "s/localhost/$registry_ip/g" ${WORKSPACE}/deploy/kustomize/contrail/1node/latest/*
fi

kubectl apply -k ${WORKSPACE}/deploy/kustomize/operator/latest/
while [[ ! $(kubectl wait crds --for=condition=Established --timeout=2m managers.contrail.juniper.net) ]]
do
  sleep 2s
done

kubectl apply -k ${WORKSPACE}/deploy/kustomize/contrail/1node/latest/
