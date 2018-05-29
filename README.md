# Kubernetes Terraform Provider

This provider is a fork of the official Kubernetes provider developed by HashiCorp.
This fork supports the following resources in addition to the official provider:

- `DaemonSets`
- `Deployments`
- `Ingress`
- `StatefulSets`

## Supported Kubernetes Versions

The latest build of this provider uses v6.0 of the kubernetes [client-go](https://github.com/kubernetes/client-go) library, and has been tested with the following Kubernetes versions:

- 1.7.x
- 1.8.x
- 1.9.x

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/sl1pm4t/terraform-provider-kubernetes`

```sh
$ mkdir -p $GOPATH/src/github.com/sl1pm4t; cd $GOPATH/src/github.com/sl1pm4t
$ git clone https://github.com/sl1pm4t/terraform-provider-kubernetes
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sl1pm4t/terraform-provider-kubernetes
$ make build
```

## Using the provider

**Provider Configuration**

##### Simplest - kubectl configuration
```hcl-terraform
provider kubernetes {
  # leave blank to pickup config from kubectl config of local system
}
```

##### Explicit configuration
```hcl-terraform
provider "kubernetes" {
  host     = "https://104.196.242.174"
  username = "ClusterMaster"
  password = "MindTheGap"

  client_certificate     = "${file("~/.kube/client-cert.pem")}"
  client_key             = "${file("~/.kube/client-key.pem")}"
  cluster_ca_certificate = "${file("~/.kube/cluster-ca-cert.pem")}"
}
```

##### Initialise provider with plugin directory

After this fork has been downloaded and built into `$GOPATH` (as in the
previous step), specify the location of the built binaries when
initialising the project:

    $ terraform init -plugin-dir=$GOPATH/bin

This step is important in order to make terraform actually use this fork
of the kubernetes provider. If you fail to do this upon the first init
the official provider is downloaded instead.
If this happens, delete the `.terraform/` folder that has been created
inside your project folder and perform the above init again.

**Deployment Resource**

```hcl-terraform
resource "kubernetes_deployment" "nginx" {

  metadata {
    name      = "nginx"
    namespace = "web"
  }

  spec {
    selector {
      app = "nginx"
    }

    template {
      metadata {
        labels {
          app = "nginx"
        }
      }

      spec {
        container {
          image = "nginx:1.8"
          name  = "app"

          resources {
            requests {
              memory = "1Gi"
              cpu    = "1"
            }

            limits {
              memory = "2Gi"
              cpu    = "2"
            }
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = "90"
            }

            initial_delay_seconds = 10
            period_seconds        = 10
          }

          liveness_probe {
            exec {
              command = ["/bin/health"]
            }

            initial_delay_seconds = 120
            period_seconds        = 15
          }

          env {
            name  = "CONFIG_FILE_LOCATION"
            value = "/etc/app/config"
          }

          port {
            container_port = 80
          }

          volume_mount {
            name       = "config"
            mount_path = "/etc/app/config"
          }
        }

        init_container {
          name  = "helloworld"
          image = "debian"
          command = ["/bin/echo", "hello", "world"]
        }

        volume {
          name = "config"

          config_map {
            name = "app-config"
          }
        }

      }
    }
  }
}
```

## Developing the Provider

### Development Environment

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-kubernetes
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
