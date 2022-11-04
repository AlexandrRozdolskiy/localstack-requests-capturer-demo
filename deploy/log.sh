#!/bin/bash

kubectl logs `kubectl get pod|grep demo|awk '{print $1}'` -c capturer