module github.com/nicholasjackson/consul-smi-controller

go 1.15

require (
	github.com/go-logr/logr v0.3.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/kr/pretty v0.2.0
	github.com/nicholasjackson/smi-controller-sdk v0.0.0-20210124163956-d67d74f3897d
	github.com/servicemeshinterface/smi-sdk-go v0.4.1
	k8s.io/client-go v0.18.8
	sigs.k8s.io/controller-runtime v0.6.0
)

replace github.com/servicemeshinterface/smi-sdk-go v0.4.1 => github.com/nicholasjackson/smi-sdk-go v0.0.0-20210123215756-d8c5233cc084
