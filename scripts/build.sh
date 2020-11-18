#!/bin/bash

operator-sdk build localhost:5000/contrail-operator:latest
docker push localhost:5000/contrail-operator:latest
