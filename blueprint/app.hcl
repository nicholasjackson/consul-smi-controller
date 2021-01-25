k8s_config "app-config" {
  depends_on = ["module.consul"]

  cluster = "k8s_cluster.dc1"
  paths = [
    "./app/consul_config.yaml",
  ]

  wait_until_ready = true
}

k8s_config "app-pods" {
  depends_on = ["k8s_config.app-config"]

  cluster = "k8s_cluster.dc1"
  paths = [
    "./app/load_test.yaml",
    "./app/web.yaml",
    "./app/apiV1.yaml",
    "./app/apiV2.yaml",
  ]

  wait_until_ready = true
}
