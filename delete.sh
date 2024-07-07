#!/bin/bash

./kubernetes/loadbalancer/delete.sh
./kubernetes/cassanadra/delete.sh
./kubernetes/token-service/delete.sh

kubectl delete all --all -n urlshortener
kubectl delete pvc --all -n urlshortener
kubectl delete pv --all -n urlshortener
kubectl delete StorageClass: --all -n urlshortener


kubectl delete namespace urlshortener