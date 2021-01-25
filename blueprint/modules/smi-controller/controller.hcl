helm "smi-controler" {
  depends_on = ["helm.cert-manager"]

  cluster = "k8s_cluster.${var.consul_k8s_cluster}"
  namespace = "smi"
  create_namespace = true

  // chart = "github.com/hashicorp/consul-helm?ref=crd-controller-base"
  chart = "github.com/nicholasjackson/smi-controller-sdk/helm//smi-controller"

  values_string = {
    "controller.enabled" = "false"
  }
}
