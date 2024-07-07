#!/bin/bash

kubectl delete -f ./kubernetes/cassandra/cluster-ip.yaml
kubectl delete -f ./kubernetes/cassandra/StatefulSet.yaml
kubectl delete -f ./kubernetes/cassandra/storageclass.yaml