helm "flagger" {
  depends_on = ["helm.cert-manager"]

  cluster = "k8s_cluster.${var.consul_k8s_cluster}"
  namespace = "smi"
  create_namespace = true

  chart = "github.com/fluxcd/flagger//charts/flagger?ref=v1.6.1"

  values_string = {
    "meshProvider" = ""
    "metricsServer" = "http://prometheus-stack-kube-prom-prometheus.default.svc:9090" 
  }
}
