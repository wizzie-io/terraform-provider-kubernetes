package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	api "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAccKubernetesClusterRoleBinding_basic(t *testing.T) {
	var conf api.ClusterRoleBinding
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	roleName := fmt.Sprintf("tf-acc-role-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_cluster_role_binding.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesClusterRoleBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterRoleBindingConfig_basic(roleName, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesClusterRoleBindingExists("kubernetes_cluster_role_binding.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "role_ref.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "role_ref.0.kind", "ClusterRole"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "role_ref.0.name", roleName),
					resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "subject.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "subject.0.kind", "Group"),
					//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.0.verbs.1", "list"),
				),
			},
			//{
			//	Config: testAccKubernetesClusterRoleBindingConfig_modified(roleName, name),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		testAccCheckKubernetesClusterRoleBindingExists("kubernetes_cluster_role_binding.test", &conf),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.#", "2"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.0.verbs.#", "3"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.0.verbs.2", "watch"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.1.api_groups.#", "1"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.1.resources.#", "1"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.1.resources.0", "deployments"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.1.verbs.#", "2"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.1.verbs.0", "get"),
			//		//resource.TestCheckResourceAttr("kubernetes_cluster_role_binding.test", "rule.1.verbs.1", "list"),
			//	),
			//},
		},
	})
}

func TestAccKubernetesClusterRoleBinding_importBasic(t *testing.T) {
	resourceName := "kubernetes_cluster_role_binding.test"
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	roleName := fmt.Sprintf("tf-acc-role-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesClusterRoleBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterRoleBindingConfig_basic(roleName, name),
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

func testAccCheckKubernetesClusterRoleBindingDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*kubernetesProvider).conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_cluster_role_binding" {
			continue
		}
		_, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.RbacV1().ClusterRoleBindings().Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Cluster Role still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesClusterRoleBindingExists(n string, obj *api.ClusterRoleBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetesProvider).conn
		_, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		out, err := conn.RbacV1().ClusterRoleBindings().Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesClusterRoleBindingConfig_basic(rolename, name string) string {
	return fmt.Sprintf(`
resource "kubernetes_cluster_role" "test" {
	metadata {
		name = "%s"
	}
	rule {
		api_groups = [""]
		resources  = ["pods", "pods/log"]
		verbs = ["get", "list"]
	}
}

resource "kubernetes_cluster_role_binding" "test" {
	metadata {
		name = "%s"
	}
	role_ref {
		name  = "%s"
		kind  = "ClusterRole"
	}
	subject {
		kind = "Group"
		name = "monitoring"
	}
}`, rolename, name, rolename)
}

//func testAccKubernetesClusterRoleBindingConfig_modified(name string) string {
//	return fmt.Sprintf(`
//resource "kubernetes_cluster_role_binding" "test" {
//	metadata {
//		name = "%s"
//	}
//	role_ref {
//		api_groups = [""]
//		resources  = ["pods", "pods/log"]
//		verbs      = ["get", "list", "watch"]
//	}
//}`, name)
//}
