#!/bin/bash

kubectl expose deployment go-k8s-example --type=LoadBalancer --name=go-k8s-example
kubectl get pod
minikube service go-k8s-example --url
