#!/bin/bash

kubectl apply -f ./kubernetes/cassandra/storageclass.yaml
kubectl apply -f ./kubernetes/cassandra/StatefulSet.yaml
kubectl apply -f ./kubernetes/cassandra/cluster-ip.yaml