module "consul" {
  source = "./modules/consul"
}

module "smi-controller" {
  source = "./modules/smi-controller"
}

module "monitoring" {
  depends_on = ["module.consul"]
  source = "./modules/monitoring"
}
