#!/bin/bash -x

kubectl delete -k deploy/kustomize/contrail/1node/latest/
kubectl delete -k deploy/kustomize/operator/latest/
kubectl delete -f deploy/crds/
kubectl delete pv  cassandra1-pv-0 zookeeper1-pv-0
sudo rm -rf \
  /mnt/cassandra \
  /mnt/zookeeper \
  /var/lib/contrail \
  /var/log/contrail \
  /var/crashes/contrail \
  /etc/cni/net.d/10-tf-cni.conf
