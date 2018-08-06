package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	api "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceKubernetesClusterRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesClusterRoleBindingCreate,
		Read:   resourceKubernetesClusterRoleBindingRead,
		Exists: resourceKubernetesClusterRoleBindingExists,
		Update: resourceKubernetesClusterRoleBindingUpdate,
		Delete: resourceKubernetesClusterRoleBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": metadataSchema("cluster role binding", true),
			"role_ref": {
				Type:        schema.TypeList,
				Description: "RoleRef can only reference a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error. See official documentation: https://v1-9.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#roleref-v1-rbac",
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: roleRefFields(),
				},
			},
			"subject": {
				Type:        schema.TypeList,
				Description: "Subjects holds references to the objects the role applies to.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: rbacSubjectFields(),
				},
			},
		},
	}
}

func resourceKubernetesClusterRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	crb := api.ClusterRoleBinding{
		ObjectMeta: metadata,
		RoleRef:    expandRoleRef(d.Get("role_ref").([]interface{})[0]),
		Subjects:   expandSubjects(d.Get("subject").([]interface{})),
	}
	log.Printf("[INFO] Creating new cluster role binding: %#v", crb)
	out, err := conn.RbacV1().ClusterRoleBindings().Create(&crb)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new cluster role binding: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesClusterRoleBindingRead(d, meta)
}

func resourceKubernetesClusterRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	_, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Reading cluster role binding %s", name)
	crb, err := conn.RbacV1().ClusterRoleBindings().Get(name, metav1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received cluster role binding: %#v", crb)
	err = d.Set("metadata", flattenMetadata(crb.ObjectMeta, d))
	if err != nil {
		return err
	}
	d.Set("role_ref", flattenRoleRef(crb.RoleRef))
	d.Set("subject", flattenSubjects(crb.Subjects))

	return nil
}

func resourceKubernetesClusterRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	_, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	crb := api.ClusterRoleBinding{
		ObjectMeta: metadata,
		RoleRef:    expandRoleRef(d.Get("role_ref").([]interface{})[0]),
		Subjects:   expandSubjects(d.Get("subject").([]interface{})),
	}

	log.Printf("[INFO] Updating cluster role binding %q: %v", name, crb)
	out, err := conn.RbacV1().ClusterRoleBindings().Update(&crb)
	if err != nil {
		return fmt.Errorf("Failed to update cluster role binding: %s", err)
	}
	log.Printf("[INFO] Submitted updated cluster role binding: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesClusterRoleBindingRead(d, meta)
}

func resourceKubernetesClusterRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	_, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deleting cluster role binding: %#v", name)
	err = conn.RbacV1().ClusterRoleBindings().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] cluster role binding %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesClusterRoleBindingExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetesProvider).conn

	_, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking cluster role binding %s", name)
	_, err = conn.RbacV1().ClusterRoleBindings().Get(name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
