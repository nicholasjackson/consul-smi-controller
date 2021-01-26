container "api_1" {
  image {
    name = "nicholasjackson/fake-service:vm-v0.13.2"
  }

  volume {
    source      = "./files/api_1.hcl"
    destination = "/config/api_1.hcl"
  }

  network { 
    name = "network.dc1"
  }

  env {
    key = "CONSUL_SERVER"
    value = "${shipyard_ip()}"
  }
  
  env {
    key = "SERVICE_ID"
    value = "api-3"
  }

  env {
    key = "LISTEN_ADDR"
    value = "0.0.0.0:9090"
  }

  env {
    key = "NAME"
    value = "API V3"
  }
  
  env {
    key = "MESSAGE"
    value = "Hello from API V3"
  }
}

