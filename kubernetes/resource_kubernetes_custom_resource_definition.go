package kubernetes

import (
	"fmt"
	"log"

	api "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"github.com/hashicorp/terraform/helper/schema"
	//api "k8s.io/api/extensions/v1beta1"

	cr "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
)

func resourceKubernetesCustomResourceDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesCustomResourceDefinitionCreate,
		Read:   resourceKubernetesCustomResourceDefinitionRead,
		Exists: resourceKubernetesCustomResourceDefinitionExists,
		Update: resourceKubernetesCustomResourceDefinitionUpdate,
		Delete: resourceKubernetesCustomResourceDefinitionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("custom resource definition", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec describes how the user wants the resources to appear",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:        schema.TypeString,
							Description: "Group is the group this resource belongs in",
							Required:    true,
							ForceNew:    true,
						},
						"name": {
							Type:        schema.TypeList,
							Description: "Group is the group this resource belongs in",
							Required:    true,
							ForceNew:    true,
							MinItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"plural": {
										Type:        schema.TypeString,
										Description: "Plural is the plural name of the resource to serve.  It must match the name of the CustomResourceDefinition-registration",
										Optional:    true,
									},
									"singular": {
										Type:        schema.TypeString,
										Description: "Singular is the singular name of the resource.  It must be all lowercase  Defaults to lowercased <kind>",
										Computed:    true,
										Optional:    true,
									},
									"short_names": {
										Type:        schema.TypeSet,
										Description: "ShortNames are short names for the resource.  It must be all lowercase.",
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
									"kind": {
										Type:        schema.TypeString,
										Description: "Kind is the serialized kind of the resource.  It is normally CamelCase and singular.",
										Required:    true,
									},
									"list_kind": {
										Type:        schema.TypeString,
										Description: "ListKind is the serialized kind of the list for this resource.  Defaults to <kind>List",
										Optional:    true,
										Computed:    true,
									},
									"categories": {
										Type:        schema.TypeSet,
										Description: "Categories is a list of grouped resources custom resources belong to (e.g. 'all')",
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
								},
							},
						},
						"scope": {
							Type:        schema.TypeString,
							Description: "Scope indicates whether this resource is cluster or namespace scoped. Default is namespaced",
							Optional:    true,
							Default:     "Namespaced",
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Version is the version this resource belongs in",
							Required:    true,
							//Deprecated:  "Use versions",
						},
						//"versions": {
						//	Type:        schema.TypeList,
						//	Description: "Version is the version this resource belongs in",
						//	Required:    true,
						//	ForceNew:    true,
						//	MinItems:    1,
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"plural": {
						//				Type:        schema.TypeString,
						//				Description: "Plural is the plural name of the resource to serve.  It must match the name of the CustomResourceDefinition-registration",
						//				Optional:    true,
						//			},
						//		},
						//	},
						//},
					},
				},
			},
		},
	}
}

func resourceKubernetesCustomResourceDefinitionCreate(d *schema.ResourceData, meta interface{}) error {
	prov := meta.(*kubernetesProvider)
	//conn := meta.(*kubernetesProvider).conn

	conn, err := api.NewForConfig(prov.cfg)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))

	crd := cr.CustomResourceDefinition{
		//TypeMeta: metav1.TypeMeta{
		//	Kind:       "CustomResourceDefinition",
		//	APIVersion: "apiextensions.k8s.io/v1beta1",
		//},
		ObjectMeta: metadata,
		Spec:       expandCustomResourceDefinitionSpec(d.Get("spec").([]interface{})),
	}
	//jCRD, _ := json.Marshal(crd)

	log.Printf("[INFO] Creating new custom resource definition: %#v", crd)

	//out := &cr.CustomResourceDefinition{}

	out, err := conn.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&crd)

	//out, err := conn.Resource(res).Create(unstr, metav1.GetOptions{}).
	//	NamespaceIfScoped(metadata.Namespace, crd.Spec.Scope == "Namespaced").
	//	Resource("customresourcedefinitions").
	//	Prefix("apis", "apiextensions").
	//	Name(metadata.Name).
	//	Body(jCRD).
	//	Do().
	//	Into(out)
	if err != nil {
		return fmt.Errorf("could not create CRD %s: %s", crd.Name, err)
	}

	log.Printf("[INFO] Submitted new custom resource definition: %#v", out)
	d.SetId(buildId(metadata))

	return resourceKubernetesCustomResourceDefinitionRead(d, meta)
}

func resourceKubernetesCustomResourceDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Reading custom resource definition %s", name)
	crd := &cr.CustomResourceDefinition{}
	err = conn.RESTClient().Get().
		Namespace(namespace).
		Resource("customresourcedefinitions").
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do().
		Into(crd)
	//res := conn.RESTClient().Get().Namespace(namespace).Name(fmt.Sprintf("customresourcedefinition/%s", name)).Do()
	if err != nil {
		return err
	}
	//err = res.Into(crd)
	//if err != nil {
	//	log.Printf("[DEBUG] Received error: %#v", err)
	//	return err
	//}
	log.Printf("[INFO] Received custom resource definition: %#v", crd)
	err = d.Set("metadata", flattenMetadata(crd.ObjectMeta, d))
	if err != nil {
		return err
	}
	d.Set("spec", crd.Spec)

	return nil
}

func resourceKubernetesCustomResourceDefinitionUpdate(d *schema.ResourceData, meta interface{}) error {
	//conn := meta.(*kubernetesProvider).conn
	//
	//namespace, name, err := idParts(d.Id())
	//if err != nil {
	//	return err
	//}
	//
	//log.Printf("[INFO] Updating custom resource definition %q: %v", name, string(data))
	//out, err := conn.CoreV1().CustomResourceDefinitions(namespace).Patch(name, pkgApi.JSONPatchType, data)
	//if err != nil {
	//	return fmt.Errorf("Failed to update custom resource definition: %s", err)
	//}
	//log.Printf("[INFO] Submitted updated custom resource definition: %#v", out)
	//d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesCustomResourceDefinitionRead(d, meta)
}

func resourceKubernetesCustomResourceDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}
	res := conn.RESTClient().Delete().Namespace(namespace).Name(fmt.Sprintf("customresourcedefinition/%s", name)).Do()
	if res.Error() != nil {
		return res.Error()
	}

	log.Printf("[INFO] Custom Resource Definition %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesCustomResourceDefinitionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetesProvider).conn

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking custom resource definition %s", name)
	res := conn.RESTClient().Get().Namespace(namespace).Name(fmt.Sprintf("customresourcedefinition/%s", name)).Do()
	if res.Error() != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	return true, err
}

func expandCustomResourceDefinitionSpec(in []interface{}) cr.CustomResourceDefinitionSpec {
	if len(in) == 0 {
		return cr.CustomResourceDefinitionSpec{}
	}
	cfgCfg := in[0].(map[string]interface{})
	crd := cr.CustomResourceDefinitionSpec{}

	if v, ok := cfgCfg["group"]; ok {
		crd.Group = v.(string)
	}
	if v, ok := cfgCfg["name"]; ok {
		crd.Names = expandCustomResourceDefinitionName(v.([]interface{}))
	}
	//if v, ok := cfgCfg["resource_names"]; ok {
	//	crd.ResourceNames = expandStringSlice(v.([]interface{}))
	//}
	//if v, ok := cfgCfg["resources"]; ok {
	//	crd.Resources = expandStringSlice(v.([]interface{}))
	//}
	//if v, ok := cfgCfg["versions"]; ok {
	//	crd.Versions = expandStringSlice(v.([]interface{}))
	//}

	return crd
}

func expandCustomResourceDefinitionName(in []interface{}) cr.CustomResourceDefinitionNames {
	n := cr.CustomResourceDefinitionNames{}

	namesCfg := in[0].(map[string]interface{})
	if v, ok := namesCfg["kind"]; ok {
		n.Kind = v.(string)
	}
	if v, ok := namesCfg["list_kind"]; ok {
		n.ListKind = v.(string)
	}
	if v, ok := namesCfg["plural"]; ok {
		n.Plural = v.(string)
	}
	if v, ok := namesCfg["singular"]; ok {
		n.Singular = v.(string)
	}

	return n
}
