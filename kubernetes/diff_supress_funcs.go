package kubernetes

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/resource"
)

func suppressEquivalentResourceQuantity(k, old, new string, d *schema.ResourceData) bool {
	oldQ, err := resource.ParseQuantity(old)
	if err != nil {
		return false
	}
	newQ, err := resource.ParseQuantity(new)
	if err != nil {
		return false
	}
	return oldQ.Cmp(newQ) == 0
}

func suppressEmptyMapValues(k, old, new string, d *schema.ResourceData) bool {
	key := strings.Replace(k, "#", new, 1)
	_, ok := d.GetOk(key)
	return !ok
}
