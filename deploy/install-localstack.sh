#!/bin/bash

brew install helm
helm repo add localstack-repo https://helm.localstack.cloud
helm list |grep localstack | helm upgrade --install localstack localstack-repo/localstack