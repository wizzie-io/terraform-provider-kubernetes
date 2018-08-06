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

func expandRoleRef(in interface{}) v1.RoleRef {
	obj := v1.RoleRef{}

	rrCfg := in.(map[string]interface{})

	if v, ok := rrCfg["api_group"]; ok {
		obj.APIGroup = v.(string)
	}
	if v, ok := rrCfg["kind"]; ok {
		obj.Kind = v.(string)
	}
	if v, ok := rrCfg["name"]; ok {
		obj.Name = v.(string)
	}

	return obj
}

func expandSubjects(in []interface{}) []v1.Subject {
	if len(in) < 1 {
		return []v1.Subject{}
	}
	subs := make([]v1.Subject, len(in))

	for i, v := range in {
		subCfg := v.(map[string]interface{})
		sub := v1.Subject{}

		if v, ok := subCfg["api_group"]; ok {
			sub.APIGroup = v.(string)
		}
		if v, ok := subCfg["kind"]; ok {
			sub.Kind = v.(string)
		}
		if v, ok := subCfg["name"]; ok {
			sub.Name = v.(string)
		}
		if v, ok := subCfg["namespace"]; ok {
			sub.Namespace = v.(string)
		}
		subs[i] = sub
	}

	return subs
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

func flattenRoleRef(in v1.RoleRef) []interface{} {
	m := make(map[string]interface{})

	m["api_group"] = in.APIGroup
	m["kind"] = in.Kind
	m["name"] = in.Name

	return []interface{}{m}
}

func flattenSubjects(in []v1.Subject) []interface{} {
	att := make([]interface{}, len(in), len(in))
	for i, n := range in {
		m := make(map[string]interface{})

		m["api_group"] = n.APIGroup
		m["kind"] = n.Kind
		m["name"] = n.Name
		m["namespace"] = n.Namespace

		att[i] = m
	}

	return att
}
