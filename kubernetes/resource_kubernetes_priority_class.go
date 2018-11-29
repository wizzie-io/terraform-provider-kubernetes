package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/api/scheduling/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
)

func resourceKubernetesPriorityClass() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesPriorityClassCreate,
		Read:   resourceKubernetesPriorityClassRead,
		Exists: resourceKubernetesPriorityClassExists,
		Update: resourceKubernetesPriorityClassUpdate,
		Delete: reosurceKubernetesPriorityClassDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": metadataSchema("priority class", true),
			"global_default": {
				Type:        schema.TypeBool,
				Description: "Indicates whether this PriorityClass should be considered as the default priority for pods that do not have any priority class",
				Optional:    true,
				Default:     false,
			},
			"value": {
				Type:        schema.TypeInt,
				Description: "The value of this priority class",
				Required:    true,
			},
		},
	}
}

func resourceKubernetesPriorityClassCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	priorityClass := v1beta1.PriorityClass{
		ObjectMeta:    metadata,
		GlobalDefault: d.Get("global_default").(bool),
		Value:         int32(d.Get("value").(int)),
	}

	log.Printf("[INFO] Creating new priority class: %#v", priorityClass)
	out, err := conn.SchedulingV1beta1().PriorityClasses().Create(&priorityClass)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new storage class: %#v", out)
	d.SetId(out.Name)
	return resourceKubernetesPriorityClassRead(d, meta)
}

func resourceKubernetesPriorityClassRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	name := d.Id()
	log.Printf("[INFO] Reading priority class %s", name)
	priorityClass, err := conn.Scheduling().PriorityClasses().Get(name, metav1.GetOptions{})

	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received priority class: %#v", priorityClass)

	err = d.Set("metadata", flattenMetadata(priorityClass.ObjectMeta, d))

	if err != nil {
		return err
	}

	d.Set("global_default", priorityClass.GlobalDefault)
	d.Set("value", priorityClass.Value)

	return nil
}

func resourceKubernetesPriorityClassUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	name := d.Id()
	ops := patchMetadata("metadata.0.", "/metadata/", d)
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}
	log.Printf("[INFO] Updating priority class %q: %v", name, string(data))
	out, err := conn.Scheduling().PriorityClasses().Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update priority class: %s", err)
	}
	log.Printf("[INFO] Submitted updated priority class: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesPriorityClassRead(d, meta)
}

func reosurceKubernetesPriorityClassDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	name := d.Id()
	log.Printf("[INFO] Deleting priority class: %#v", name)
	err := conn.Scheduling().PriorityClasses().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Priority class %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesPriorityClassExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetesProvider).conn

	name := d.Id()
	log.Printf("[INFO] checking storage class %s", name)
	_, err := conn.Scheduling().PriorityClasses().Get(name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok &&
			statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
