helm "cert-manager" {
  depends_on = ["k8s_config.cert-manager"]

  cluster = "k8s_cluster.dc1"

  chart = "github.com/jetstack/cert-manager//deploy/charts/cert-manager?ref=v1.1.0"
  
  values_string = {
    "installCRDs" = "true"
    "image.tag" = "v1.1.0"
    "cainjector.image.tag" = "v1.1.0"
    "webhook.image.tag" = "v1.1.0"
  }
 
  health_check {
    timeout = "60s"
    pods = ["app=webhook"]
  }
}

k8s_config "cert-manager" {
  cluster = "k8s_cluster.dc1"
  paths = [
    "./k8sconfig/cert-manager.crds.yaml",
  ]

  wait_until_ready = true
}
