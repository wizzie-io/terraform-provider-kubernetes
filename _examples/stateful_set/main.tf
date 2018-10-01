resource "kubernetes_stateful_set" "test" {
  metadata {
    name = "test"
  }

  spec {
    replicas = 2

    selector {
      app = "httpbin"
    }

    service_name = "test"

    update_strategy {
      type = "RollingUpdate"

      rolling_update {
        partition = 1
      }
    }

    template {
      metadata {
        labels {
          app = "httpbin"
        }
      }

      spec {
        container {
          image = "citizenstig/httpbin"
          name  = "httpbin"

          volume_mount {
            name       = "pvc"
            mount_path = "/data"
          }
        }
      }
    }

    volume_claim_templates {
      metadata {
        name = "pvc"
      }

      spec {
        access_modes = ["ReadWriteOnce"]

        resources {
          requests {
            storage = "2Gi"
          }
        }
      }
    }
  }
}
