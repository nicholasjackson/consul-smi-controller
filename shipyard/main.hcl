variable "consul_k8s_cluster" {
  default = "dc1"
}

variable "consul_k8s_network" {
  default = "dc1"
}

variable "monitoring_k8s_cluster" {
  default = "dc1"
}

variable "consul_smi_controller_enabled" {
  default = true
}

variable "smi_controller_k8s_cluster" {
  default = "dc1"
}

variable "smi_controller_k8s_network" {
  default = "dc1"
}

variable "smi_controller_enabled" {
  default = false
}

variable "smi_controller_webhook_enabled" {
  default = false
}

variable "smi_controller_webhook_port" {
  default = 9443
}

variable "smi_controller_namespace" {
  default = "shipyard"
}

variable "smi_controller_additional_dns" {
  default = "smi-webhook.shipyard.svc"
}

variable "consul_acls_enabled" {
  default = true
}

variable "consul_tls_enabled" {
  default = true
}

network "dc1" {
  subnet = "10.5.0.0/16"
}

k8s_cluster "dc1" {
  driver = "k3s"

  nodes = 1

  network {
    name = "network.dc1"
  }
}

output "KUBECONFIG" {
  value = k8s_config("dc1")
}

module "consul" {
  source = "github.com/shipyard-run/blueprints//modules/kubernetes-consul?ref=a6ee46a479de9fad68af54a90081b2c577d66195"
}

# Create an ingress which exposes the locally running webhook from kubernetes
ingress "smi-webhook" {
  source {
    driver = "k8s"

    config {
      cluster = "k8s_cluster.dc1"
      port    = 9443
    }
  }

  destination {
    driver = "local"

    config {
      address = "localhost"
      port    = 9443
    }

  }
}