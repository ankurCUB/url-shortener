#!/bin/bash

kubectl create namespace urlshortener

./kubernetes/token-service/deployment.sh
./kubernetes/cassandra/deployment.sh
./kubernetes/loadbalancer/deployment.sh