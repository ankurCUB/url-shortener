#!/bin/bash

kubectl delete -f ./kubernetes/token-service/tokenservice.yaml
kubectl delete -f ./kubernetes/token-service/mysqlpv.yaml