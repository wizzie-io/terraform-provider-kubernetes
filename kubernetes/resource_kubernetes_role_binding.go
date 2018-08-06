package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	api "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceKubernetesRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesRoleBindingCreate,
		Read:   resourceKubernetesRoleBindingRead,
		Exists: resourceKubernetesRoleBindingExists,
		Update: resourceKubernetesRoleBindingUpdate,
		Delete: resourceKubernetesRoleBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("role binding", false),
			"role_ref": {
				Type:        schema.TypeList,
				Description: "RoleRef can reference a Role in the current namespace or a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error.",
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

func resourceKubernetesRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	rb := api.RoleBinding{
		ObjectMeta: metadata,
		RoleRef:    expandRoleRef(d.Get("role_ref").([]interface{})[0]),
		Subjects:   expandSubjects(d.Get("subject").([]interface{})),
	}
	log.Printf("[INFO] Creating new role binding: %#v", rb)
	out, err := conn.RbacV1().RoleBindings(metadata.Namespace).Create(&rb)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new role binding: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesRoleBindingRead(d, meta)
}

func resourceKubernetesRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Reading role binding %s", name)
	crb, err := conn.RbacV1().RoleBindings(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received role binding: %#v", crb)
	err = d.Set("metadata", flattenMetadata(crb.ObjectMeta, d))
	if err != nil {
		return err
	}
	d.Set("role_ref", flattenRoleRef(crb.RoleRef))
	d.Set("subject", flattenSubjects(crb.Subjects))

	return nil
}

func resourceKubernetesRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	crb := api.RoleBinding{
		ObjectMeta: metadata,
		RoleRef:    expandRoleRef(d.Get("role_ref").([]interface{})[0]),
		Subjects:   expandSubjects(d.Get("subject").([]interface{})),
	}

	log.Printf("[INFO] Updating role binding %q: %v", name, crb)
	out, err := conn.RbacV1().RoleBindings(namespace).Update(&crb)
	if err != nil {
		return fmt.Errorf("Failed to update role binding: %s", err)
	}
	log.Printf("[INFO] Submitted updated role binding: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesRoleBindingRead(d, meta)
}

func resourceKubernetesRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deleting role binding: %#v", name)
	err = conn.RbacV1().RoleBindings(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] role binding %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesRoleBindingExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking role binding %s", name)
	_, err = conn.RbacV1().RoleBindings(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
