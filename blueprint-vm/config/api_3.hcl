service {
  name = "api"
  id = "api-3"
  port = 9090
  tags = ["v3"]

  connect { 
    sidecar_service {
      proxy {
      }
    }
  }
}
