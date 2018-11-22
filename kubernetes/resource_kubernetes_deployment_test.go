package kubernetes

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	appsv1 "k8s.io/api/apps/v1"
)

func TestAccKubernetesDeployment_minimal(t *testing.T) {
	t.Parallel()

	var conf appsv1.Deployment
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_deployment.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentConfig_minimal(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.image", "nginx:1.7.8"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.name", "tf-acc-test"),
				),
			},
		},
	})
}

func TestAccKubernetesDeployment_basic(t *testing.T) {
	t.Parallel()

	var conf appsv1.Deployment
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_deployment.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentConfig_basic(name, 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.TestAnnotationTwo", "two"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one", "TestAnnotationTwo": "two"}),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelThree", "three"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelTwo": "two", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.image", "nginx:1.7.8"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.name", "tf-acc-test"),
				),
			},
			{
				Config: testAccKubernetesDeploymentConfig_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.Different", "1234"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one", "Different": "1234"}),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelThree", "three"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.image", "nginx:1.7.9"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.name", "tf-acc-test"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.paused", "true"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.progress_deadline_seconds", "30"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.revision_history_limit", "4"),
				),
			},
			{
				Config: testAccKubernetesDeploymentConfig_basic(name, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.annotations.TestAnnotationTwo", "two"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one", "TestAnnotationTwo": "two"}),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.TestLabelThree", "three"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelTwo": "two", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_deployment.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.image", "nginx:1.7.8"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.name", "tf-acc-test"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.paused", "false"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.progress_deadline_seconds", "600"),
				),
			},
		},
	})
}

func TestAccKubernetesDeployment_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "kubernetes_deployment.test"
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentConfig_basic(name, 2),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
		},
	})
}

func TestAccKubernetesDeployment_with_template_metadata(t *testing.T) {
	t.Parallel()

	var conf appsv1.Deployment

	depName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentConfigWithTemplateMetadata(depName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.image", imageName),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.annotations.prometheus.io/scrape", "true"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.annotations.prometheus.io/scheme", "https"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.annotations.prometheus.io/port", "4000"),
					pause(),
				),
			},
			{
				Config: testAccKubernetesDeploymentConfigWithTemplateMetadataModified(depName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.container.0.image", imageName),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.labels.foo", "bar"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.annotations.prometheus.io/scrape", "true"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.annotations.prometheus.io/scheme", "http"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.metadata.0.annotations.prometheus.io/port", "8080"),
					pause(),
				),
			},
		},
	})
}

func TestAccKubernetesDeployment_initContainer(t *testing.T) {
	t.Parallel()

	var conf appsv1.Deployment
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_deployment.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentWithInitContainer(name, "nginx:1.7.8"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.template.0.spec.0.init_container.0.image", "alpine"),
				),
			},
		},
	})
}
func TestAccKubernetesDeployment_noTopLevelLabels(t *testing.T) {
	t.Parallel()

	var conf appsv1.Deployment
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_deployment.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentWithNoTopLevelLabels(name, "nginx:1.7.8"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "metadata.0.labels.%", "0"),
				),
			},
		},
	})
}

func TestAccKubernetesDeployment_strategy(t *testing.T) {
	t.Parallel()

	var conf appsv1.Deployment
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_deployment.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesDeploymentWithStrategy(name, "Recreate"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.type", "Recreate"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.#", "0"),
				),
			},
			{
				Config: testAccKubernetesDeploymentWithStrategy(name, "RollingUpdate"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.type", "RollingUpdate"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.#", "1"),
				),
			},
			{
				Config: testAccKubernetesDeploymentWithStrategy(name, "Recreate"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.type", "Recreate"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.#", "0"),
				),
			},
			{
				Config: testAccKubernetesDeploymentWithRollingUpdateStrategy(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.type", "RollingUpdate"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.0.max_surge", "50%"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.0.max_surge", "50%"),
				),
			},
			{
				Config: testAccKubernetesDeploymentWithRollingUpdateStrategy2(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesDeploymentExists("kubernetes_deployment.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.type", "RollingUpdate"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.0.max_surge", "25%"),
					resource.TestCheckResourceAttr("kubernetes_deployment.test", "spec.0.strategy.0.rolling_update.0.max_surge", "25%"),
				),
			},
		},
	})
}

func pause() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(1 * time.Minute)
		return nil
	}
}

func testAccCheckKubernetesDeploymentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*kubernetesProvider)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_deployment" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := readDeployment(conn, namespace, name)
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Deployment still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesDeploymentExists(n string, obj *appsv1.Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetesProvider)

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		out, err := readDeployment(conn, namespace, name)

		*obj = *out
		return nil
	}
}

func testAccKubernetesDeploymentConfig_minimal(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
  }
  spec {
		replicas = 3
    selector {
      foo = "bar"
    }
    template {
			metadata {
				labels {
					foo = "bar"
				}
			}
			spec {
				container {
					image = "nginx:1.7.8"
					name  = "tf-acc-test"
				}
			}
    }
  }
}
`, name)
}

func testAccKubernetesDeploymentConfig_basic(name string, replicas int) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    annotations {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    replicas = %d

    selector {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    template {
      metadata {
        labels {
          TestLabelOne   = "one"
          TestLabelTwo   = "two"
          TestLabelThree = "three"
        }
      }

      spec {
        container {
          image = "nginx:1.7.8"
          name  = "tf-acc-test"
        }
      }
    }
  }
}
`, name, replicas)
}

func testAccKubernetesDeploymentConfig_modified(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    annotations {
      TestAnnotationOne = "one"
      Different         = "1234"
    }

    labels {
      TestLabelOne   = "one"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    paused                    = true
    progress_deadline_seconds = 30
    revision_history_limit    = 4

    selector {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    template {
      metadata {
        labels {
          TestLabelOne   = "one"
          TestLabelTwo   = "two"
          TestLabelThree = "three"
        }
      }

      spec {
        container {
          image = "nginx:1.7.9"
          name  = "tf-acc-test"
        }
      }
    }
  }
}`, name)
}

func testAccKubernetesDeploymentConfigWithTemplateMetadata(depName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
    labels {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector {
      foo = "bar"
      Test = "TfAcceptanceTest"
    }
    template {
		metadata {
			labels {
				foo = "bar"
				Test = "TfAcceptanceTest"
			}
			annotations {
				"prometheus.io/scrape" = "true"
				"prometheus.io/scheme" = "https"
				"prometheus.io/port"   = "4000"
                "kubernetes.io/egress-bandwidth" = "1M"
			}
		}
		spec {
			container {
				image = "%s"
				name  = "containername"
			}
		}
    }
  }
}
`, depName, imageName)
}

func testAccKubernetesDeploymentConfigWithTemplateMetadataModified(depName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
    labels {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector {
			foo = "bar"
      Test = "TfAcceptanceTest"
		}
    template {
			metadata {
				labels {
					foo = "bar"
					Test = "TfAcceptanceTest"
				}
				annotations {
					"prometheus.io/scrape" = "true"
					"prometheus.io/scheme" = "http"
					"prometheus.io/port"   = "8080"
				}
			}
			spec {
				container {
					image = "%s"
					name  = "containername"
				}
			}
    }
  }
}
`, depName, imageName)
}

func testAccKubernetesDeploymentWithInitContainer(depName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
    labels {
			foo = "bar"
    }
  }

  spec {
    selector {
			foo = "bar"
		}
    template {
			metadata {
				labels {
					foo = "bar"
				}
			}
			spec {
				init_container {
					name = "hello"
					image = "alpine"
					command = ["echo", "'hello'"]
				}
				container {
					image = "%s"
					name  = "containername"
				}
			}
    }
  }
}
`, depName, imageName)
}

func testAccKubernetesDeploymentWithNoTopLevelLabels(depName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
  }

  spec {
    selector {
			foo = "bar"
		}
    template {
			metadata {
				labels {
					foo = "bar"
				}
			}
			spec {
				container {
					image = "%s"
					name  = "containername"
				}
			}
    }
  }
}
`, depName, imageName)
}

func testAccKubernetesDeploymentWithStrategy(depName, strategy string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
  }

  spec {
    selector {
      foo = "bar"
    }

    strategy {
      type = "%s"
    }

    template {
      metadata {
        labels {
          foo = "bar"
        }
      }

      spec {
        container {
          image = "alpine"
          name  = "containername"
        }
      }
    }
  }
}
`, depName, strategy)
}

func testAccKubernetesDeploymentWithRollingUpdateStrategy(depName string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
  }

  spec {
    selector {
      foo = "bar"
    }

    strategy {
      type = "RollingUpdate"
      rolling_update {
        max_surge       = "50%%"
        max_unavailable = "50%%"
      }
    }

    template {
      metadata {
        labels {
          foo = "bar"
        }
      }

      spec {
        container {
          image = "alpine"
          name  = "containername"
        }
      }
    }
  }
}
`, depName)
}

func testAccKubernetesDeploymentWithRollingUpdateStrategy2(depName string) string {
	return fmt.Sprintf(`
resource "kubernetes_deployment" "test" {
  metadata {
    name = "%s"
  }

  spec {
    selector {
      foo = "bar"
    }

    strategy {
      type = "RollingUpdate"
      rolling_update {}
    }

    template {
      metadata {
        labels {
          foo = "bar"
        }
      }

      spec {
        container {
          image = "alpine"
          name  = "containername"
        }
      }
    }
  }
}
`, depName)
}
