resource "kubernetes_deployment" "httpbin" {
  metadata {
    name = "httpbin"
  }

  spec {
    selector {
      foo = "bar"
      app = "amaze"
    }

    template {
      metadata {
        labels {
          foo = "bar"
          app = "amaze"
        }

        annotations {
          "prometheus.io/scrape" = "true"
        }
      }

      spec {
        container {
          image = "citizenstig/httpbin"
          name  = "app"

          port {
            container_port = "8000"
          }

          readiness_probe {
            http_get {
              path = "/healthz"
              port = "8000"
            }
          }
        }
      }
    }
  }
}
