//
// Install Consul using the helm chart.
//
helm "consul" {
  cluster = "k8s_cluster.${var.consul_k8s_cluster}"

  // chart = "github.com/hashicorp/consul-helm?ref=crd-controller-base"
  chart = "github.com/hashicorp/consul-helm?ref=v0.28.0"
  values = "./helm/consul-values.yaml"

  health_check {
    timeout = "60s"
    pods = ["app=consul"]
  }
}

ingress "consul" {
  source {
    driver = "local"
    
    config {
      port = 8500
    }
  }
  
  destination {
    driver = "k8s"
    
    config {
      cluster = "k8s_cluster.${var.consul_k8s_cluster}"
      address = "consul-server.default.svc"
      port = 8500
    }
  }
}

ingress "consul-rpc" {
  source {
    driver = "local"
    
    config {
      port = 8300
    }
  }
  
  destination {
    driver = "k8s"
    
    config {
      cluster = "k8s_cluster.${var.consul_k8s_cluster}"
      address = "consul-server.default.svc"
      port = 8300
    }
  }
}

ingress "consul-lan-serf" {
  source {
    driver = "local"
    
    config {
      port = 8301
    }
  }
  
  destination {
    driver = "k8s"
    
    config {
      cluster = "k8s_cluster.${var.consul_k8s_cluster}"
      address = "consul-server.default.svc"
      port = 8301
    }
  }
}
