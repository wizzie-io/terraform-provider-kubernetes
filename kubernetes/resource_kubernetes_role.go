package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	api "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceKubernetesRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesRoleCreate,
		Read:   resourceKubernetesRoleRead,
		Exists: resourceKubernetesRoleExists,
		Update: resourceKubernetesRoleUpdate,
		Delete: resourceKubernetesRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("role", true),
			"rule": {
				Type:        schema.TypeList,
				Description: "List of PolicyRules for this Role",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: policyRuleFields(),
				},
			},
		},
	}
}

func resourceKubernetesRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	cRole := api.Role{
		ObjectMeta: metadata,
		Rules:      expandClusterRoleRule(d.Get("rule").([]interface{})),
	}
	log.Printf("[INFO] Creating new role: %#v", cRole)
	out, err := conn.RbacV1().Roles(metadata.Namespace).Create(&cRole)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new role: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesRoleRead(d, meta)
}

func resourceKubernetesRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Reading role %s", name)
	cRole, err := conn.RbacV1().Roles(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received role: %#v", cRole)
	err = d.Set("metadata", flattenMetadata(cRole.ObjectMeta, d))
	if err != nil {
		return err
	}
	d.Set("rule", flattenClusterRoleRules(cRole.Rules))

	return nil
}

func resourceKubernetesRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	cRole := api.Role{
		ObjectMeta: metadata,
		Rules:      expandClusterRoleRule(d.Get("rule").([]interface{})),
	}

	log.Printf("[INFO] Updating role %q: %v", name, cRole)
	out, err := conn.RbacV1().Roles(namespace).Update(&cRole)
	if err != nil {
		return fmt.Errorf("Failed to update role: %s", err)
	}
	log.Printf("[INFO] Submitted updated role: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesRoleRead(d, meta)
}

func resourceKubernetesRoleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deleting role: %#v", name)
	err = conn.RbacV1().Roles(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] role %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking role %s", name)
	_, err = conn.RbacV1().Roles(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
