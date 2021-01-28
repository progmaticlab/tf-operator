#!/bin/bash -e
my_file="$(readlink -e "$0")"
my_dir="$(dirname $my_file)"
WORKSPACE=${WORKSPACE:-$HOME/tf-operator}

export CONTAINER_REGISTRY=${CONTAINER_REGISTRY:-"tungstenfabric"}
export CONTRAIL_CONTAINER_TAG=${CONTRAIL_CONTAINER_TAG:-"latest"}

count=$(echo "$CONTROLLER_NODES $AGENT_NODES" | tr " " "\n" | sort -u |  tr "\n" " " | awk -F ' ' '{print NF}' )
if [[ $count > 1 ]] && [[ $CONTAINER_REGISTRY == localhost* ]] ; then
  # multinode setup with registry on one node
  registry_ip=$(hostname -i)
  export CONTAINER_REGISTRY=$(echo $CONTAINER_REGISTRY | sed "s/loclahost/$registry_ip/g")
fi

OPERATOR_KUSTOMIZATION_FILE="$my_dir/../deploy/kustomize/operator/latest/kustomization.yaml"
"$my_dir/jinja2_render.py" < ${OPERATOR_KUSTOMIZATION_FILE}.j2 > ${OPERATOR_KUSTOMIZATION_FILE}
kubectl apply -k $my_dir/../deploy/kustomize/operator/latest/
while [[ ! $(kubectl wait crds --for=condition=Established --timeout=2m managers.contrail.juniper.net) ]]
do
  sleep 2s
done

CONTRAIL_KUSTOMIZATION_FILE="$my_dir/../deploy/kustomize/contrail/1node/latest/kustomization.yaml"
"$my_dir/jinja2_render.py" < ${CONTRAIL_KUSTOMIZATION_FILE}.j2 > ${CONTRAIL_KUSTOMIZATION_FILE}
kubectl apply -k $my_dir/../deploy/kustomize/contrail/1node/latest/
