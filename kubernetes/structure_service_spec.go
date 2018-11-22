package kubernetes

import (
	gversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/version"
)

// Flatteners

func flattenIntOrString(in intstr.IntOrString) int {
	return in.IntValue()
}

func flattenServicePort(in []v1.ServicePort) []interface{} {
	att := make([]interface{}, len(in), len(in))
	for i, n := range in {
		m := make(map[string]interface{})
		m["name"] = n.Name
		m["protocol"] = string(n.Protocol)
		m["port"] = int(n.Port)
		m["target_port"] = flattenIntOrString(n.TargetPort)
		m["node_port"] = int(n.NodePort)

		att[i] = m
	}
	return att
}

func flattenServiceSpec(in v1.ServiceSpec) []interface{} {
	att := make(map[string]interface{})
	if in.ClusterIP != "" {
		att["cluster_ip"] = in.ClusterIP
	}
	if len(in.ExternalIPs) > 0 {
		att["external_ips"] = newStringSet(schema.HashString, in.ExternalIPs)
	}
	if in.ExternalName != "" {
		att["external_name"] = in.ExternalName
	}
	if in.ExternalTrafficPolicy != "" {
		att["external_traffic_policy"] = in.ExternalTrafficPolicy
	}
	if in.HealthCheckNodePort > 0 {
		att["health_check_node_port"] = in.HealthCheckNodePort
	}
	if in.LoadBalancerIP != "" {
		att["load_balancer_ip"] = in.LoadBalancerIP
	}
	if len(in.LoadBalancerSourceRanges) > 0 {
		att["load_balancer_source_ranges"] = newStringSet(schema.HashString, in.LoadBalancerSourceRanges)
	}
	if len(in.Ports) > 0 {
		att["port"] = flattenServicePort(in.Ports)
	}

	att["publish_not_ready_addresses"] = in.PublishNotReadyAddresses

	if len(in.Selector) > 0 {
		att["selector"] = in.Selector
	}
	if in.SessionAffinity != "" {
		att["session_affinity"] = string(in.SessionAffinity)
	}
	if in.SessionAffinityConfig != nil && in.SessionAffinityConfig.ClientIP != nil {
		att["session_affinity_config"] = flattenSessionAffinityConfig(in.SessionAffinityConfig)
	}
	if in.Type != "" {
		att["type"] = string(in.Type)
	}
	return []interface{}{att}
}

func flattenLoadBalancerIngress(in []v1.LoadBalancerIngress) []interface{} {
	out := make([]interface{}, len(in), len(in))
	for i, ingress := range in {
		att := make(map[string]interface{})

		att["ip"] = ingress.IP
		att["hostname"] = ingress.Hostname

		out[i] = att
	}
	return out
}

func flattenSessionAffinityConfig(in *v1.SessionAffinityConfig) []interface{} {
	if in == nil {
		return nil
	}
	att := make(map[string]interface{})
	if in.ClientIP != nil {
		clientIPAtt := make(map[string]interface{})
		if in.ClientIP.TimeoutSeconds != nil && *in.ClientIP.TimeoutSeconds != 10800 {
			clientIPAtt["timeout_seconds"] = int(*in.ClientIP.TimeoutSeconds)
			att["client_ip_config"] = []interface{}{clientIPAtt}
			return []interface{}{att}
		}
	}
	return nil
}

// Expanders

func expandIntOrString(in int) intstr.IntOrString {
	return intstr.FromInt(in)
}

func expandServicePort(l []interface{}) []v1.ServicePort {
	if len(l) == 0 || l[0] == nil {
		return []v1.ServicePort{}
	}
	obj := make([]v1.ServicePort, len(l), len(l))
	for i, n := range l {
		cfg := n.(map[string]interface{})
		obj[i] = v1.ServicePort{
			Port:       int32(cfg["port"].(int)),
			TargetPort: expandIntOrString(cfg["target_port"].(int)),
		}
		if v, ok := cfg["name"].(string); ok {
			obj[i].Name = v
		}
		if v, ok := cfg["protocol"].(string); ok {
			obj[i].Protocol = v1.Protocol(v)
		}
		if v, ok := cfg["node_port"].(int); ok {
			obj[i].NodePort = int32(v)
		}
	}
	return obj
}

func expandServiceSpec(l []interface{}) v1.ServiceSpec {
	if len(l) == 0 || l[0] == nil {
		return v1.ServiceSpec{}
	}
	in := l[0].(map[string]interface{})
	obj := v1.ServiceSpec{}

	// process type first, as it's needed for conditional handling of other attributes
	if v, ok := in["type"].(string); ok {
		obj.Type = v1.ServiceType(v)
	}

	if v, ok := in["cluster_ip"].(string); ok {
		obj.ClusterIP = v
	}
	if v, ok := in["external_ips"].(*schema.Set); ok && v.Len() > 0 {
		obj.ExternalIPs = sliceOfString(v.List())
	}
	if v, ok := in["external_name"].(string); ok {
		obj.ExternalName = v
	}
	if v, ok := in["external_traffic_policy"].(string); ok && (obj.Type == v1.ServiceTypeNodePort || obj.Type == v1.ServiceTypeLoadBalancer) {
		obj.ExternalTrafficPolicy = v1.ServiceExternalTrafficPolicyType(v)
	}
	if v, ok := in["health_check_node_port"].(int); ok {
		obj.HealthCheckNodePort = int32(v)
	}
	if v, ok := in["load_balancer_ip"].(string); ok {
		obj.LoadBalancerIP = v
	}
	if v, ok := in["load_balancer_source_ranges"].(*schema.Set); ok && v.Len() > 0 {
		obj.LoadBalancerSourceRanges = sliceOfString(v.List())
	}
	if v, ok := in["port"].([]interface{}); ok && len(v) > 0 {
		obj.Ports = expandServicePort(v)
	}
	if v, ok := in["publish_not_ready_addresses"].(bool); ok {
		obj.PublishNotReadyAddresses = v
	}
	if v, ok := in["selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Selector = expandStringMap(v)
	}
	if v, ok := in["session_affinity"].(string); ok {
		obj.SessionAffinity = v1.ServiceAffinity(v)
	}
	if v, ok := in["session_affinity_config"].([]interface{}); ok && len(v) > 0 {
		obj.SessionAffinityConfig = expandSessionAffinityConfig(v)
	}
	return obj
}

func expandSessionAffinityConfig(l []interface{}) *v1.SessionAffinityConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	var obj *v1.SessionAffinityConfig

	for _, n := range l {
		cfg := n.(map[string]interface{})
		if v, ok := cfg["client_ip_config"].([]interface{}); ok && len(v) > 0 {
			for _, n2 := range v {
				cfg2 := n2.(map[string]interface{})
				if v2, ok2 := cfg2["timeout_seconds"].(int); ok2 {
					obj = &v1.SessionAffinityConfig{
						ClientIP: &v1.ClientIPConfig{
							TimeoutSeconds: ptrToInt32(int32(v2)),
						},
					}
				}
			}
		}

	}
	return obj
}

// Patch Ops

func patchServiceSpec(keyPrefix, pathPrefix string, d *schema.ResourceData, v *version.Info) (PatchOperations, error) {
	ops := make([]PatchOperation, 0, 0)
	if d.HasChange(keyPrefix + "selector") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "selector",
			Value: d.Get(keyPrefix + "selector").(map[string]interface{}),
		})
	}
	if d.HasChange(keyPrefix + "type") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "type",
			Value: d.Get(keyPrefix + "type").(string),
		})
	}
	if d.HasChange(keyPrefix + "session_affinity") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "sessionAffinity",
			Value: d.Get(keyPrefix + "session_affinity").(string),
		})
	}
	if d.HasChange(keyPrefix + "load_balancer_ip") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "loadBalancerIP",
			Value: d.Get(keyPrefix + "load_balancer_ip").(string),
		})
	}
	if d.HasChange(keyPrefix + "external_traffic_policy") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "externalTrafficPolicy",
			Value: d.Get(keyPrefix + "external_traffic_policy").(string),
		})
	}
	if d.HasChange(keyPrefix + "load_balancer_source_ranges") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "loadBalancerSourceRanges",
			Value: d.Get(keyPrefix + "load_balancer_source_ranges").(*schema.Set).List(),
		})
	}
	if d.HasChange(keyPrefix + "port") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "ports",
			Value: expandServicePort(d.Get(keyPrefix + "port").([]interface{})),
		})
	}
	if d.HasChange(keyPrefix + "external_ips") {
		k8sVersion, err := gversion.NewVersion(v.String())
		if err != nil {
			return nil, err
		}
		v1_7_0, _ := gversion.NewVersion("1.7.0")
		if k8sVersion.LessThan(v1_7_0) {
			// If we haven't done this the deprecated field would have priority
			ops = append(ops, &ReplaceOperation{
				Path:  pathPrefix + "deprecatedPublicIPs",
				Value: nil,
			})
		}

		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "externalIPs",
			Value: d.Get(keyPrefix + "external_ips").(*schema.Set).List(),
		})
	}
	if d.HasChange(keyPrefix + "external_name") {
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "externalName",
			Value: d.Get(keyPrefix + "external_name").(string),
		})
	}
	return ops, nil
}
