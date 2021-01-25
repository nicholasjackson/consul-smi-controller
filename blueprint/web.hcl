ingress "web" {
  source {
    driver = "local"
    
    config {
      port = 18080
    }
  }
  
  destination {
    driver = "k8s"
    
    config {
      cluster = "k8s_cluster.${var.consul_k8s_cluster}"
      address = "web-service.default.svc"
      port = 9090
    }
  }
}
