#!/bin/bash -e

type go >/dev/null 2>&1 || {
  export PATH=$PATH:/usr/local/go/bin
}

export  CGO_ENABLED=1
operator-sdk build  localhost:5000/tf-operator:latest
sudo docker push localhost:5000/tf-operator:latest

# build CRDS container
sudo docker build --tag localhost:5000/tf-crdsloader:latest deploy/crds
sudo docker push localhost:5000/tf-crdsloader:latest
