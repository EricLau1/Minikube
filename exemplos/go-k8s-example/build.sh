#!/bin/bash

docker build -t go-k8s-example:1.0 .
kubectl create -f ./k8s/deployment.yaml
