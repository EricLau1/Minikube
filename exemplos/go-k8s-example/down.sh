#!/bin/bash

kubectl delete service go-k8s-example
kubectl get service
kubectl delete deployment go-k8s-example
kubectl get deployment
kubectl delete pod go-k8s-example
kubectl get pod
