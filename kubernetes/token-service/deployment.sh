#!/bin/bash

kubectl apply -f ./kubernetes/token-service/mysqlpv.yaml
kubectl apply -f ./kubernetes/token-service/tokenservice.yaml