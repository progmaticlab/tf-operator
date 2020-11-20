#!/bin/bash

operator-sdk build localhost:5000/tf-operator:latest
sudo docker push localhost:5000/tf-operator:latest

# build CRDS container
sudo docker build --tag localhost:5000/tf-crdsloader:latest deploy/crds
sudo docker push localhost:5000/tf-crdsloader:latest
