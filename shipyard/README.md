---
title: Consul Service Mesh on Kubernetes with Monitoring
author: Nic Jackson
slug: k8s_consul_stack
shipyard_version: ">= 0.2.11"
---

# Consul Service Mesh on Kubernetes with Monitoring


This blueprint creates a Kubernetes cluster and installs the following elements:

* Consul Service Mesh With CRDs
* Prometheus
* Loki
* Grafana
* Flagger
* Cert Manager
* SMI Controller for Consul
* Example Application

To access Grafana the following details can be used:

* user: admin
* pass: admin

ACLs are disabled for Consul
