#!/bin/bash
kubectl apply -f consul_config.yaml
kubectl apply -f web.yaml -f apiV1.yaml -f apiV2.yaml
