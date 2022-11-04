#!/bin/bash

kubectl apply -f deploy/rbac.yaml
kubectl delete deploy demo 
kubectl apply -f deploy/deploy.yaml
