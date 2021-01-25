k8s_config "grafana_secret" {
  cluster = "k8s_cluster.${var.monitoring_k8s_cluster}"

  paths = [
    "./k8sconfig/grafana_secret.yaml",
  ]

  wait_until_ready = true
}

helm "grafana" {
  cluster = "k8s_cluster.${var.monitoring_k8s_cluster}"

  chart = "github.com/grafana/helm-charts/charts//grafana"
  values = "./helm/grafana_values.yaml"
}

ingress "grafana" {
  source {
    driver = "local"
    
    config {
      port = 8080
    }
  }
  
  destination {
    driver = "k8s"
    
    config {
      cluster = "k8s_cluster.${var.monitoring_k8s_cluster}"
      address = "grafana.default.svc"
      port = 80
    }
  }
}
