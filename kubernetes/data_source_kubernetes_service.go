package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceKubernetesService() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceKubernetesService().Schema)

	addRequiredFieldsToSchema(dsSchema, "metadata")
	addRequiredFieldsToSchema(dsSchema["metadata"].Elem.(*schema.Resource).Schema, "name")
	addRequiredFieldsToSchema(dsSchema["metadata"].Elem.(*schema.Resource).Schema, "namespace")

	return &schema.Resource{
		Read: dataSourceKubernetesServiceRead,

		Schema: dsSchema,
	}
}

func dataSourceKubernetesServiceRead(d *schema.ResourceData, meta interface{}) error {
	om := meta_v1.ObjectMeta{
		Namespace: d.Get("metadata.0.namespace").(string),
		Name:      d.Get("metadata.0.name").(string),
	}
	d.SetId(buildId(om))

	return resourceKubernetesServiceRead(d, meta)
}
