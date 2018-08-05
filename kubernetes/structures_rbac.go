package kubernetes

import (
	"k8s.io/api/rbac/v1"
)

func expandClusterRoleRule(in []interface{}) []v1.PolicyRule {
	if len(in) == 0 {
		return []v1.PolicyRule{}
	}
	rules := make([]v1.PolicyRule, len(in))

	for i, rule := range in {
		r := v1.PolicyRule{}

		ruleCfg := rule.(map[string]interface{})
		if v, ok := ruleCfg["api_groups"]; ok {
			r.APIGroups = expandStringSlice(v.([]interface{}))
		}
		if v, ok := ruleCfg["non_resource_urls"]; ok {
			r.NonResourceURLs = expandStringSlice(v.([]interface{}))
		}
		if v, ok := ruleCfg["resource_names"]; ok {
			r.ResourceNames = expandStringSlice(v.([]interface{}))
		}
		if v, ok := ruleCfg["resources"]; ok {
			r.Resources = expandStringSlice(v.([]interface{}))
		}
		if v, ok := ruleCfg["verbs"]; ok {
			r.Verbs = expandStringSlice(v.([]interface{}))
		}

		rules[i] = r
	}

	return rules
}

// Flatteners
func flattenClusterRoleRules(in []v1.PolicyRule) []interface{} {
	att := make([]interface{}, len(in), len(in))
	for i, n := range in {
		m := make(map[string]interface{})

		m["api_groups"] = n.APIGroups
		m["non_resource_urls"] = n.NonResourceURLs
		m["resource_names"] = n.ResourceNames
		m["resources"] = n.Resources
		m["verbs"] = n.Verbs

		att[i] = m
	}

	return att
}
